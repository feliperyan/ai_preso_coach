import 'package:ai_preso_01/camera_view.dart';
import 'package:ai_preso_01/model.dart';
import 'package:flutter/material.dart';
import 'package:ai_preso_01/welcome_view.dart';
import 'package:provider/provider.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();

  runApp(ChangeNotifierProvider(
      create: (context) => TheModel(),
      child: MaterialApp(
          title: 'CamApp',
          theme: ThemeData(
            useMaterial3: true,
            colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
          ),
          home: const FirstRoute())));
}

class FirstRoute extends StatelessWidget {
  const FirstRoute({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<TheModel>(builder: thingy);
  }
}

Widget thingy(BuildContext context, TheModel model, Widget? child) {
  return Scaffold(
    appBar: AppBar(
      title: const Text('First Route'),
    ),
    body: Center(
        child: Column(
      children: [
        Welcome(image: model.getPicture()),
        ElevatedButton(
          child: const Text('Open route'),
          onPressed: () {
            // Navigate to second route when tapped.
            Navigator.push(context,
                MaterialPageRoute(builder: (context) => const CameraView()));
          },
        ),
      ],
    )),
  );
}
