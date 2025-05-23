import 'package:TrackMini/services/api_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter_config/flutter_config.dart';
import 'package:provider/provider.dart';
import './events.dart';
import './profile.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await FlutterConfig.loadEnvVariables();

  runApp(
    MultiProvider(
      providers: [
        Provider(
          create:
              (context) => ApiService(baseUrl: FlutterConfig.get('BASE_URL')),
        ),
        ChangeNotifierProvider(create: (context)),
      ],
    ),
  );
}

class TrackMini extends StatelessWidget {
  const TrackMini({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'TrackMini',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.grey),
      ),
      home: const Home(title: 'TrackMini'),
    );
  }
}

class Home extends StatefulWidget {
  const Home({super.key, required this.title});

  final String title;

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  int _selectedIndex = 0;

  final List<Widget> _screensList = <Widget>[
    const EventsScreen(),
    const ProfileScreen(),
  ];

  void _onNavigate(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(child: _screensList.elementAt(_selectedIndex)),
      bottomNavigationBar: BottomNavigationBar(
        items: const <BottomNavigationBarItem>[
          BottomNavigationBarItem(icon: Icon(Icons.event), label: 'Events'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
        currentIndex: _selectedIndex,
        selectedItemColor: Colors.blue,
        onTap: _onNavigate,
      ),
    );
  }
}
