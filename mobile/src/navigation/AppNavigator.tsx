import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';

// Placeholder screens - to be implemented
const PlaceholderScreen = ({ route }: any) => {
  return null;
};

const Stack = createNativeStackNavigator();
const Tab = createBottomTabNavigator();

function MainTabs() {
  return (
    <Tab.Navigator>
      <Tab.Screen name="Feed" component={PlaceholderScreen} />
      <Tab.Screen name="Search" component={PlaceholderScreen} />
      <Tab.Screen name="Upload" component={PlaceholderScreen} />
      <Tab.Screen name="Notifications" component={PlaceholderScreen} />
      <Tab.Screen name="Profile" component={PlaceholderScreen} />
    </Tab.Navigator>
  );
}

export default function AppNavigator() {
  return (
    <NavigationContainer>
      <Stack.Navigator>
        <Stack.Screen
          name="Main"
          component={MainTabs}
          options={{ headerShown: false }}
        />
        <Stack.Screen name="Video" component={PlaceholderScreen} />
        <Stack.Screen name="UserProfile" component={PlaceholderScreen} />
        <Stack.Screen name="Login" component={PlaceholderScreen} />
        <Stack.Screen name="Signup" component={PlaceholderScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}