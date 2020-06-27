<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no" name="viewport">
    <meta content="yes" name="apple-mobile-web-app-capable">
    <meta content="black" name="apple-mobile-web-app-status-bar-style">
    <meta content="telephone=no" name="format-detection">
    <meta content="email=no" name="format-detection">
    <title>支付成功</title>
    <script>
    (function(designWidth, maxWidth) {
        var doc = document,
        win = window,
        docEl = doc.documentElement,
        remStyle = document.createElement("style"),
        tid;

        function refreshRem() {
            var width = docEl.getBoundingClientRect().width;
            maxWidth = maxWidth || 540;
            width>maxWidth && (width=maxWidth);
            var rem = width * 100 / designWidth;
            remStyle.innerHTML = 'html{font-size:' + rem + 'px;}';
        }

        if (docEl.firstElementChild) {
            docEl.firstElementChild.appendChild(remStyle);
        } else {
            var wrap = doc.createElement("div");
            wrap.appendChild(remStyle);
            doc.write(wrap.innerHTML);
            wrap = null;
        }
        //要等 wiewport 设置好后才能执行 refreshRem，不然 refreshRem 会执行2次；
        refreshRem();

        win.addEventListener("resize", function() {
            clearTimeout(tid); //防止执行两次
            tid = setTimeout(refreshRem, 300);
        }, false);

        win.addEventListener("pageshow", function(e) {
            if (e.persisted) { // 浏览器后退的时候重新计算
                clearTimeout(tid);
                tid = setTimeout(refreshRem, 300);
            }
        }, false);

        if (doc.readyState === "complete") {
            doc.body.style.fontSize = "16px";
        } else {
            doc.addEventListener("DOMContentLoaded", function(e) {
                doc.body.style.fontSize = "16px";
            }, false);
        }
    })(1080,1080);
</script>
    <style>
        body,img{
            margin: 0;
            padding: 0;
        }
        img{
            width: 4.8rem;
	        height: 5.8rem;
            position: absolute;
            top: 35%;
            left: 50%;
            transform: translate(-50%,-50%); 
        }
    </style>
</head>
<body>
    <img src="/static/success.png" alt="">
</body>
</html>