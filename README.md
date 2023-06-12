## SNS-connections

## About
My second personal project that I've been writing in Go with Gin framework. I just wanted to build a simple microservice first but then come up with idea to connect different Social Networks such as Instagram, TikTok, YouTube, etc and make one full-fledged web application. The main point is to share videos between those services, for example I have an Instagram account, created content and I also want some of these content to be in TikTok as well, I can just reupload them one by one, but it will take quite a lot of time. Here when my project come in handy. The sharing may not only be limited to videos but also posts, images, even live stories...

## Features
Simple registration&login. But in order to share videos between platforms you need to connect those services to your user profile, for YouTube you need to connect Google account for Instagram you need Facebook account. Videos firstly uploaded to Drive then imported to designated social network.

## Tiktok?
TikTok API requires my web-app to be hosted somewhere, right now I'm not planning to do this. I created **Docker** container though. In the future maybe I will create domain too, who knows.

## Tech
I used Gin framework for backend with simple local database SQLite. I'm not using any frontend frameworks since I just wanted to practise and improve my Go skills. Though, I added some simple HTML, CSS(Bootstrap) and also pure JavaScript for functionality, but they are in horrible state right now(just look at the code), so I plan to improve them and make pages halfway decent at least.

## Screenshots

![Home page](https://i.imgur.com/I8MkqGb.png)
![Videos page_1](https://i.imgur.com/zFMG0L2.png)
![Videos page_2](https://i.imgur.com/XZcVKWt.png)
![Videos page_3](https://i.imgur.com/Ks5clZu.png)
![Login page](https://i.imgur.com/wdQIoPz.png)
![Registration page](https://i.imgur.com/ncrKrOW.png)

## Demo videos(download links)
[Connecting google account and retrieving videos from your youtube channel](https://cdn-141.anonfiles.com/L0lav1w9zf/cc6fbc1b-1686566383/connect+youtube.mp4)

[Connecting facebook account and retrieving videos from your instagram account](https://cdn-149.anonfiles.com/Oaldv6wbz4/82f57892-1686566513/connect+instagram.mp4)