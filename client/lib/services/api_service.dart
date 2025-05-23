import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  final String baseUrl;

  ApiService({required this.baseUrl});

  Future<dynamic> get(String endpoint) async {
    final Uri uri = Uri.parse('$baseUrl/$endpoint');
    try {
      final response = await http.get(uri);
      return _handleResponse(response);
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  Future<dynamic> post(String endpoint, dynamic data) async {
    final Uri uri = Uri.parse('$baseUrl/$endpoint');
    final headers = {'Content-Type': 'application/json'};
    final String body = jsonEncode(data);
    try {
      final response = await http.post(uri, headers: headers, body: body);
      return _handleResponse(response);
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  Future<dynamic> put(String endpoint, dynamic data) async {
    final Uri uri = Uri.parse('$baseUrl/$endpoint');
    final headers = {'Content-Type': 'application/json'};
    final String body = jsonEncode(data);
    try {
      final response = await http.put(uri, headers: headers, body: body);
      return _handleResponse(response);
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  Future<dynamic> delete(String endpoint) async {
    final Uri uri = Uri.parse('$baseUrl/$endpoint');
    try {
      final response = await http.delete(uri);
      return _handleResponse(response);
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  dynamic _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (response.body.isNotEmpty) {
        return jsonDecode(response.body);
      }
      return null;
    } else {
      throw Exception('API Error: ${response.statusCode} - ${response.body}');
    }
  }
}
