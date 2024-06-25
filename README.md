
# run app
make run app="example"

ffmpeg -i video.mp4 -vf scale=1280:-1 -c:v h264 -c:a aac \
  -hls_time 10 -hls_playlist_type vod \
  -hls_segment_filename 'http://localhost:8086/origin/ewogICAgImZzIjogIm9zIiwKICAgICJvc0NvbmZpZyI6IHsKICAgICAgICAiYmFzZVBhdGgiOiAiZGF0YSIKICAgIH0KfQ==/02/segments/segment_%03d.ts' -hls_segment_type fmp4 \
  -master_pl_name 'index.m3u8' \
  -f hls 'http://localhost:8086/api/example/origin/02/manifest'
