import 'dart:convert';

import 'package:ai_preso_01/model.dart';
import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:http/http.dart' as http;

// CameraApp is the Main Application.
class CameraView extends StatefulWidget {
  // Default Constructor
  const CameraView({super.key});

  @override
  State<CameraView> createState() => _CameraViewState();
}

class _CameraViewState extends State<CameraView> {
  late CameraController controller;
  bool isInited = false;

  List<dynamic> _data = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();

    List<CameraDescription> cameras;
    availableCameras().then((List<CameraDescription> onValue) {
      cameras = onValue;
      controller = CameraController(cameras[0], ResolutionPreset.medium);

      controller.initialize().then((_) {
        if (!mounted) {
          return;
        }
        setState(() {});
        isInited = true;
      }).catchError((Object e) {
        if (e is CameraException) {
          switch (e.code) {
            case 'CameraAccessDenied':
              // Handle access errors here.
              break;
            default:
              // Handle other errors here.
              break;
          }
        }
      });
    });
  }

  Future<XFile?> takePicture() async {
    if (!controller.value.isInitialized) {
      print('cant take picture, controller not inited');
      return null;
    }
    try {
      final XFile file = await controller.takePicture();
      return file;
    } on CameraException catch (e) {
      print('error on controller.takePicture() was: ${e.code}');
      return null;
    }
  }

  @override
  void dispose() {
    // print('disposing...');
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!isInited) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }

    return Scaffold(
        appBar: AppBar(title: const Text('Camera View')),
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(
              width: 300,
              height: 260,
              child: CameraPreview(controller),
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                onTakePictureButtonPress(context);
              },
              // onPressed: () => {},
              child: const Text('Capture'),
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                fetchData();
              },
              // onPressed: () => {},
              child: const Text('API Call'),
            )
          ],
        )));
  }

  void onTakePictureButtonPress(BuildContext context) {
    takePicture().then((XFile? file) {
      print('onTakePictureButtonPress triggered');
      var model = context.read<TheModel>();
      model.setPicture(file!);
    });
    return;
  }

  Future<void> fetchData() async {
    print('API call');
    final response =
        await http.get(Uri.parse('https://jsonplaceholder.typicode.com/posts'));
    if (response.statusCode == 200) {
      setState(() {
        _data = json.decode(response.body);
        _isLoading = false;
        print(_data);
      });
    } else {
      throw Exception('Failed to load data');
    }
  }

  //   if (mounted) {
  //     setState(() {
  //       picture = file;
  //     });

  //     print('file saved to ${file?.path}');
  //     file?.saveTo('/home/felipe').then((_) {
  //       file.length().then((val) {
  //         print('size: ${val / 0.001}');
  //       });
  //     });
  //   }
  // });
}
