import cv2
from flask import Flask, Response

app = Flask(__name__)

# Initialize the camera
cap = cv2.VideoCapture(0)  # 0 is the default camera

if not cap.isOpened():
    cap = cv2.VideoCapture(1)

if not cap.isOpened():
    cap = cv2.VideoCapture(2)

if not cap.isOpened():
    cap = cv2.VideoCapture(3)

def generate_frames():
    boundary = "--123456789000000000000987654321"
    while True:
        success, frame = cap.read()
        if not success:
            break
        else:
            # Encode the frame in JPEG format
            ret, buffer = cv2.imencode('.jpg', frame)
            frame = buffer.tobytes()

            # Construct the response as a byte-stream
            yield (b'--boundary\r\n'
                   b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')

@app.route('/')
def video_feed():
    # Video streaming route.
    return Response(generate_frames(),
                    mimetype='multipart/x-mixed-replace; boundary=boundary')

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=7000, debug=True)
