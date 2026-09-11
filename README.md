ffmpeg -i https://s3.ap-south-1.amazonaws.com/hsl.abhinayjangde.dev/videos/sample.mp4 -codec:v libx264 -codec:a aac -preset medium -crf 23 -start_number 0 -hls_time 6 -hls_list_size 0 -f hls output.m3u8



# Common LLD Case Studies & Problems

- Easy Level: Tic-Tac-Toe, Vending Machine, Parking Lot, and Logger Framework.- Medium Level: LRU Cache, Elevator System, Rate Limiter, and Movie Ticket Booking (BookMyShow).
- Advanced Level: Splitwise, Chess Game, Ride-Sharing (Uber), and Food Delivery (Zomato/Swiggy).