# buzzer-beater
Raspberry Pi powered internet of doors.

## What is buzzer-beater

I had a package stolen and decided that I'd install a surveillance camera. After
sleeping on that decision I realized that it was an awful patch. The problem isn't
that I can't see the people on my porch -- it's that I'm not at my porch to accept
packages. But I can be! The ultimate goal of buzzer-beater is to always be at my
door. I should be able to close it when I need to, see who's there, talk to them,
and open it if I trust them.

## Why not an existing product

There are commercial smart door bells but I want to build my own. In some ways
the existing products don't serve my exact needs but it really isn't lack of
existing product that is pushing me to build a smart door myself. It is just
wanting to do it.

## What will it do

Eventually I hope to (in order of appearance)
 1. Display a message to door ringers
 2. Notify door owner of ringers
 3. Send photo of ringer to users' phones
 4. Allow users to select message to send to ringers
 5. Allow lock/unlock of door with positive physical confirmation

## Current Functional Roadmap
  - [x] Users local cli creation and basic authentication over server
  - [x] User path web interface
  - [ ] Message model post/patch/delete web interface
  - [ ] Ring model with plain text logging
  - [ ] Post Ring on button press
  - [ ] Display Message on Ring post

## Feature Roadmap
 - [ ] Button press generates message on door
 - [ ] Button press sends notification to phone
 - [ ] Button press sends photo to phone
 - [ ] Send message from phone to door
 - [ ] Unlock/Lock door from phone
