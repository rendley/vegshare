  local_cam:
    # Эта настройка заставляет mediamtx забирать поток с камеры
    # вместо того, чтобы ждать его публикации через ffmpeg.
    source: rtsp://192.168.0.60:554/user=admin&password=zz123456&channel=1&stream=0.sdp
    # Указываем TCP-транспорт, как в вашей команде ffmpeg, для большей надежности
    rtspTransport: tcp



paths:
  all:
    # Конфигурация для всех путей (пока пустая)

  local_cam:
    # Запускаем ffmpeg при старте mediamtx.
    # Эта команда забирает поток с камеры, полностью его перекодирует
    # для максимальной совместимости (как в вашем рабочем примере)
    # и публикует его в этот же самый путь (local_cam) по RTMP.
    runOnInit: >-
      ffmpeg -rtsp_transport tcp
      -i "rtsp://192.168.0.60:554/user=admin&password=zz123456&channel=1&stream=0.sdp"
      -c:v libx264 -profile:v baseline -pix_fmt yuv420p -bsf:v h264_mp4toannexb
      -c:a copy
      -f flv rtmp://localhost:1935/local_cam
    # Перезапускать ffmpeg, если он упадет.
    runOnInitRestart: yes
