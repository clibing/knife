package server

const HTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, minimum-scale=1.0,maximum-scale=1.0,user-scalable=no">
    <title>上传文件</title>
    <style>
        body {
            margin: 15px auto auto 20px;
        }

        input[type='text'],
        input[type='button'] {
            font-size: 16px;
            color: #000;
            border-radius: 4px;
            border: 1px solid #ccc;
        }

        .inputDiv {
            margin: 10px auto;
        }

        .dropzone {
            width: 300px;
            height: 100px;
            border: 1px dashed gray;
            text-align: center;
            padding: 40px;
            font-size: 20px;
            color: darkgray;
        }

        .inputDiv span {
            float: left;
            width: 40px;
        }

        .inputDiv.result span {
            width: 20px;
        }

        .inputDiv input[id='token'] {
            width: 280px;
        }

        .inputDiv input.verify {
            width: 50px;
        }

        .inputDiv input[id='path'] {
            width: 335px;
        }

        .inputDiv input[id='save'] {
            float: right;
            margin-right: 5px;
        }

        .inputDiv .sub {
            width: 380px;
            height: 30px;
        }
		.dangerous {
			color: red;
		}
		.progress {
            width: 0px;
            height: 0px;
            line-height: 0px;
            border-bottom: 2px solid green;
            display: none;
		}
    </style>
</head>

<body>
    <div class="progress" id="progress"></div>
    <form method="post" enctype="multipart/form-data" id="upload">
        <div class="inputDiv dangerous" id="errMsg"></div>
		{{ if .token }} 
			<div class="inputDiv">
				<span>凭证: </span>
				<input type="text" id="token" value="" name="token" placeholder="输入token" />
				<input type="button" value="验证" class="verify" onclick="verify()" />
			</div>
		{{ end }}
        <div class="inputDiv">
            <span>目录: </span><input type="text" id="path" name="path" value="" placeholder="上传文件的目录, 不能以/开头或者结尾" />
        </div>
        <div class="inputDiv">
            <div>
                <span>文件: </span><input type="file" id="file" name="file" />
            </div>
        </div>
        <div class="dropzone" ondragover="event.preventDefault();" ondrop="handleDrop(event)">将文件拖到此处或点击选择文件</div>

        <div class="inputDiv">
            <div class="sub">
                <input type="submit" id="save" name="save" />
            </div>
        </div>

        <div class="inputDiv result">
            <ol class="filelist">
            </ol>
        </div>
    </form>

    <script>
		var waitClear = false;
		setInterval(function(){
			var errMsg = document.getElementById("errMsg").innerHTML;
			console.log("定时检查错误: ", errMsg);
			if (errMsg !== "" && !waitClear) {
				waitClear = true;
				setTimeout(function(){
					console.log("清理错误提示", errMsg);
					document.getElementById("errMsg").innerHTML = "";
				}, 5000);
			};
		}, 1000);

        function verify() {
            var token = document.getElementById("token").value;
            if (token == "" || path == "") {
				document.getElementById("errMsg").innerHTML = "请填写凭证";
                return;
            }
            var xhr = new XMLHttpRequest();
            xhr.open("GET", "/verify?token=" + token);
            xhr.onload = function () {
                console.log(xhr.responseText);
                if (xhr.status === 200) {
                    var result = JSON.parse(xhr.responseText)
                    if (result.code != 200) {
                        console.log("验证失败！");
						document.getElementById("errMsg").innerHTML = xhr.responseText;
                        return
                    }
                    document.getElementById("token").value = result.data;
					document.getElementById("errMsg").innerHTML = "";
                } else {
                    console.log("验证失败: " + xhr.responseText);
					document.getElementById("errMsg").innerHTML = xhr.responseText;
                }
            };
            xhr.send()
        }

        function handleDrop(e) {
            e.stopPropagation();
            e.preventDefault();

            var files = e.dataTransfer.files; //获取被拖动的文件列表

            var formData = new FormData();
			{{ if .token }} 
				var v =  document.getElementById("token").value;
				if (v === "") {
					document.getElementById("errMsg").innerHTML = "凭证为空";
					return
				}
				formData.set("token", v);
			{{ end }} 
            formData.set("path", document.getElementById("path").value);

            for (var i = 0; i < files.length; i++) {
                var file = files[i];
                console.log(file);
                formData.append('file', file); //将文件添加到FormData对象中
            }
            formData.append('ajax', 1); // 是否为 ajax请求组

            var xhr = new XMLHttpRequest();

            xhr.upload.addEventListener("progress", updateProgress);
            xhr.upload.addEventListener("load" , transferComplete);
            xhr.upload.addEventListener("error", transferFailed  );
            xhr.upload.addEventListener("abort", transferCanceled);

            xhr.open("POST", "/upload"); //设置上传URL为服务器接收文件的地址

            // 上传进度
            function updateProgress (event) {
                // 如果 lengthComputable 属性的值是 false，那么意味着总字节数是未知并且 total 的值为零。
                if (event.lengthComputable) {
                    let p = event.loaded / event.total * 100;
                    console.log('上传进度：' + p + '%') // 一个百分比进度
                    let v = p * 380 / 100;
                    if (v > 380) {
                        v = 380
                    }
                    console.log("width value: "+ v)
                    document.getElementById('progress').style.width = (v) + "px";
                } else {
                    // 总大小未知时不能计算进度信息
                }
            }
            function transferComplete(event) {
                console.log("上传完成");
                document.getElementById('progress').style.width = "380px";
                setTimeout(function(){
                    document.getElementById('progress').style.display= "none";
				}, 3000);
            }
            function transferFailed(event) {
                console.log("上传失败");
                document.getElementById('progress').style.width = "0px";
            }

            function transferCanceled(event) {
                console.log("取消上传");
                document.getElementById('progress').style.width = "0px";
            }

            xhr.onload = function () {
                if (xhr.status === 200) {
                    console.log("文件上传成功！");
                    var result = JSON.parse(xhr.responseText)
                    if (result.code === 200) {
                        var names = result.data;
                        for (var i = 0; i < names.length; i++) {
                            addToFileList(names[i]); //显示已上传的文件名称
                        }
						{{ if .token }} 
							document.getElementById("token").value="";
						{{ end }} 
                        return
                    }

                } else {
					console.log("文件上传失败！错误信息: " + xhr.responseText);
					document.getElementById("errMsg").innerHTML = xhr.responseText;

                }
            };
            xhr.send(formData);
            document.getElementById('progress').style.display= "block"; 
        }

        function addToFileList(fileName) {
            var ulElement = document.querySelector('.filelist');
            var liElement = document.createElement('li');
            liElement.textContent = fileName;
            ulElement.appendChild(liElement);
        }
    </script>
</body>

</html>
`
