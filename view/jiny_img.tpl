<html>
	<head>
		<title>扫描二维码付款</title>
		<meta http-equiv="content-type" content="text/html;charset=utf-8">
	</head>
<body>
<style>
body{
	background-color:#000;
}
#wrapper {
	
	width:100%;
	  
}

.content {
    width:300px;
	height:350px;
	margin-left:40%;
	margin-top:15%;
	background-color:#fff;
}
.foot{
	margin-left:30px;
}
</style> 

<div id="wrapper">  
    <div id="cell">
        <div class="content">
		<p><img id="code" src="{{.img_path}}"></p>
		<p class="foot">请使用{{.title}}扫描以上的二维码</p>
		
		</div>
    </div>
</div>	

</body>
</html>