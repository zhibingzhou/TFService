<!doctype html>
<html>
  	<head>
	</head>
	<body>
		<form id="postfrm" method="{{.Api_create_method}}" action="{{.Api_create_url}}">
			{{range $k, $v := .Api_form_param}}
				<input type="hidden" name="{{$k}}" value="{{$v}}">
			{{end}}
		</form>
正在跳转到支付界面，请稍后.....
	</body>
<script>

function htmlencode(s){
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(s));
    return div.innerHTML;
}
function htmldecode(s){
    var div = document.createElement('div');
    div.innerHTML = s;
    return div.innerText || div.textContent;
}
document.getElementById("postfrm").action=htmldecode("{{.Api_create_url}}");

document.getElementById("postfrm").submit();
</script>
</html>