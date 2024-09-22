import 'dart:typed_data';

import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

class Welcome extends StatelessWidget {
  final XFile? image;
  const Welcome({super.key, this.image});

  @override
  Widget build(BuildContext context) {
    if (image == null) {
      return const Text('No image selected.');
    } else {
      return FutureBuilder<Uint8List>(
        future: image!.readAsBytes(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done &&
              snapshot.hasData) {
            return Image.memory(snapshot.data!);
          } else if (snapshot.hasError) {
            return const Text('Error loading image.');
          } else {
            return const CircularProgressIndicator();
          }
        },
      );
    }
  }
}
