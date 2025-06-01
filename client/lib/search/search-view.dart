import 'package:flutter/material.dart';

class SearchView extends StatelessWidget {
  const SearchView({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Search Page')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              // Input Field 1: Keywords
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Keywords',
                  hintText: 'e.g., "Software Engineer", "Marketing"',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.search),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter keywords';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16.0),

              // Input Field 2: Location
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Location',
                  hintText: 'e.g., "New York, NY", "Remote"',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.location_on),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a location';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16.0),

              // Input Field 3: Category/Industry
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Category/Industry',
                  hintText: 'e.g., "Technology", "Healthcare"',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.category),
                ),
              ),
              const SizedBox(height: 16.0),

              // Input Field 4: Salary Range (Example of a numeric input)
              TextFormField(
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Minimum Salary',
                  hintText: 'e.g., "50000"',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.attach_money),
                ),
                validator: (value) {
                  if (value != null &&
                      value.isNotEmpty &&
                      int.tryParse(value) == null) {
                    return 'Please enter a valid number';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 24.0),

              // Search Button
              ElevatedButton(
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Processing Search...')),
                  );
                },
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12.0),
                ),
                child: const Text('Search', style: TextStyle(fontSize: 18.0)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
