import './api_service.dart';
import 'package:flutter_config/flutter_config.dart';

enum EventType {
  predefinedTournaments(1),
  openTournaments(2),
  teamTournaments(3),
  freestyleTournaments(4),
  seasonTournaments(5);

  const EventType(this.value);
  final int value;
}

class Event {
  final String id;
  final String name;
  final EventType type;

  Event({required this.id, required this.name, required this.type});

  factory Event.fromJson(Map<String, dynamic> json) {
    return Event(id: json['id'], name: json['name'], type: json['type']);
  }
}

class EventsService {}
