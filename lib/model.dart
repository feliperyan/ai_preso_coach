import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

class TheModel extends ChangeNotifier {
  final List<XFile> _pictures = [];

  void setPicture(XFile file) {
    _pictures.add(file);
    notifyListeners();
  }

  void clearPicture() {
    _pictures.clear();
    notifyListeners();
  }

  XFile? getPicture() {
    if (_pictures.isNotEmpty) {
      return _pictures[0];
    }

    return null;
  }
}
