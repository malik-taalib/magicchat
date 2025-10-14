# Magic Chat Mobile App

React Native mobile application for Magic Chat - TikTok-style video platform.

## Tech Stack

- **React Native** with **Expo**
- **TypeScript**
- **React Navigation** (Stack + Bottom Tabs)
- **Zustand** (State Management)
- **Axios** (API Client)
- **Expo Video** (Video Player)
- **Expo Camera** (Video Recording)
- **Expo Notifications** (Push Notifications)

## Project Structure

```
mobile/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── VideoPlayer/
│   │   ├── VideoCard/
│   │   └── shared/
│   ├── screens/          # Screen components (to be implemented)
│   ├── navigation/       # Navigation configuration
│   ├── services/         # API client and services
│   │   └── api/
│   ├── store/           # Zustand stores
│   ├── hooks/           # Custom React hooks
│   └── types/           # TypeScript types and interfaces
├── App.tsx              # App entry point
└── package.json
```

## Setup

1. **Install dependencies:**
   ```bash
   cd mobile
   npm install
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   Update `EXPO_PUBLIC_API_URL` with your backend URL.

3. **Run the app:**
   ```bash
   # iOS
   npm run ios

   # Android
   npm run android

   # Web (for testing)
   npm run web
   ```

## Development

### Running on Device

1. Install Expo Go app on your phone
2. Run `npm start`
3. Scan QR code with Expo Go (Android) or Camera (iOS)

### Building for Production

```bash
# Install EAS CLI
npm install -g eas-cli

# Configure EAS
eas login
eas build:configure

# Build for iOS
eas build --platform ios

# Build for Android
eas build --platform android
```

## API Integration

The app is configured to work with the Magic Chat Go backend. API endpoints are defined in `src/services/api/`:

- **auth.ts** - Authentication (login, register, logout)
- **videos.ts** - Video feed, likes, comments
- **users.ts** - User profiles, follow/unfollow

## Features to Implement

- [ ] Authentication screens (Login/Signup)
- [ ] For You feed with vertical scroll
- [ ] Following feed
- [ ] Video player with controls
- [ ] Like, comment, share interactions
- [ ] Video upload with camera
- [ ] User profile pages
- [ ] Search functionality
- [ ] Push notifications
- [ ] Offline support

## Notes

- All API calls use JWT authentication via AsyncStorage
- Video player uses Expo Video
- Navigation uses React Navigation v6
- State management with Zustand
- TypeScript for type safety