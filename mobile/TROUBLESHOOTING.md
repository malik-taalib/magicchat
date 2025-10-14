# Troubleshooting Guide - Magic Chat Mobile

## Starting the iOS Simulator

### Option 1: Using npm scripts
```bash
cd mobile
npm run ios
```

### Option 2: Manual start
```bash
cd mobile
npm start
# Then press 'i' to open iOS simulator
```

### Option 3: Start with tunnel (if on different network)
```bash
cd mobile
npx expo start --tunnel
```

## Common Issues & Fixes

### Issue 1: Port 8081 already in use
**Solution:**
```bash
# Kill the process using port 8081
lsof -ti:8081 | xargs kill -9

# Or use a different port
npx expo start --port 8082
```

### Issue 2: "Command not found: expo"
**Solution:**
```bash
npm install -g expo-cli
# OR use npx
npx expo start
```

### Issue 3: iOS Simulator not opening
**Solution:**
```bash
# Make sure Xcode is installed
xcode-select --install

# Open simulator manually
open -a Simulator

# Then in the Expo terminal, press 'i'
```

### Issue 4: "Unable to resolve module"
**Solution:**
```bash
# Clear cache and reinstall
rm -rf node_modules
npm install
npx expo start --clear
```

### Issue 5: Metro bundler errors
**Solution:**
```bash
# Clear Metro cache
npx expo start --clear

# Or completely reset
watchman watch-del-all
rm -rf node_modules
npm install
npx expo start --clear
```

### Issue 6: TypeScript errors
**Solution:**
```bash
# Regenerate expo types
npx expo customize tsconfig.json
```

## Checking Prerequisites

### Verify installations:
```bash
# Check Node version (should be 18+)
node --version

# Check npm version
npm --version

# Check if Xcode is installed
xcodebuild -version

# Check available iOS simulators
xcrun simctl list devices
```

## Running on Physical Device

### Using Expo Go (Easiest)
1. Install "Expo Go" from App Store on your iPhone
2. Run `npm start` in mobile folder
3. Scan QR code with Camera app (iOS) or Expo Go (Android)
4. Make sure phone and computer are on same WiFi

### Using Tunnel (if WiFi doesn't work)
```bash
npm start -- --tunnel
```

## Environment Setup

Make sure you have a `.env` file:
```bash
cd mobile
cp .env.example .env
```

Edit `.env` and set your API URL:
```
EXPO_PUBLIC_API_URL=http://localhost:8080
```

**Note for iOS Simulator:** Use your computer's IP address instead of localhost:
```
EXPO_PUBLIC_API_URL=http://192.168.1.XXX:8080
```

To find your IP:
```bash
# macOS
ipconfig getifaddr en0

# or
ifconfig | grep "inet " | grep -v 127.0.0.1
```

## Current Known Issues

✅ **Fixed Issues:**
- Added `react-native-safe-area-context`
- Added `react-native-screens`
- Fixed placeholder screens (were returning null)
- Added SafeAreaProvider wrapper
- Fixed styles for placeholder screens

## Quick Start Checklist

- [ ] Node.js installed (v18+)
- [ ] Xcode installed (for iOS)
- [ ] Navigate to mobile folder: `cd mobile`
- [ ] Dependencies installed: `npm install`
- [ ] Environment file created: `cp .env.example .env`
- [ ] No processes on port 8081: `lsof -ti:8081`
- [ ] Run: `npm start` or `npm run ios`

## Getting Help

If you're still having issues:
1. Check the Expo terminal output for specific errors
2. Look at the iOS simulator console logs
3. Try running with verbose logging: `npx expo start --verbose`
4. Check React Native docs: https://reactnative.dev/
5. Check Expo docs: https://docs.expo.dev/