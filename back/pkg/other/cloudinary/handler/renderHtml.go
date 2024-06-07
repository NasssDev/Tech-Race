package handler

import (
	"github.com/cloudinarace/entity"
	"html/template"
	"net/http"
)

func DisplayVideoHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Display Video</title>
</head>
<body>
    <h1>Uploaded Video</h1>
    <video width="320" height="240" controls>
        <source src="{{.Src}}">
    </video>

<button id="upload_widget" class="cloudinary-button">Upload files</button>

<script src="https://upload-widget.cloudinary.com/latest/global/all.js" type="text/javascript"></script>  

<script type="text/javascript">  
var myWidget = cloudinary.createUploadWidget({
  cloudName: '{{ .CloudName }}', 
  uploadPreset: 'my_preset'}, (error, result) => { 
    if (!error && result && result.event === "success") { 
      console.log('Done! Here is the image info: ', result.info); 
    }
  }
)

document.getElementById("upload_widget").addEventListener("click", function(){
    myWidget.open();
  }, false);
</script>

</body>
</html>
`
		t, err := template.New("display").Parse(tmpl)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		data := struct {
			Src       string
			CloudName string
		}{
			Src:       "https://res.cloudinary.com/demo/video/upload/f_auto/q_auto/samples/cld-sample-video.mp4",
			CloudName: entity.CloudName,
		}
		t.Execute(w, data)
	}
}

func DisplayImageHandler(entity *entity.ContextEntity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Display Image</title>
</head>
<body>
    <h1>Transformed Image</h1>
    <img src="https://res.cloudinary.com/{{ .CloudNameSRC }}" alt="Transformed Image">

<button id="upload_widget" class="cloudinary-button">Upload files</button>

<script src="https://upload-widget.cloudinary.com/latest/global/all.js" type="text/javascript"></script>  

<script type="text/javascript">  
var myWidget = cloudinary.createUploadWidget({
  cloudName: '{{ .CloudName }}', 
  uploadPreset: 'my_preset'}, (error, result) => { 
    if (!error && result && result.event === "success") { 
      console.log('Done! Here is the image info: ', result.info); 
    }
  }
)

document.getElementById("upload_widget").addEventListener("click", function(){
    myWidget.open();
  }, false);
</script>
</body>
</html>
`
		t, err := template.New("display").Parse(tmpl)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		data := struct {
			CloudNameSRC string
			CloudName    string
		}{
			CloudNameSRC: entity.CloudName + "/image/upload/r_max/e_sepia/quickstart_butterfly",
			CloudName:    entity.CloudName,
		}

		t.Execute(w, data)
	}
}
