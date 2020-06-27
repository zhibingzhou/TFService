<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no" name="viewport">
    <meta content="yes" name="apple-mobile-web-app-capable">
    <meta content="black" name="apple-mobile-web-app-status-bar-style">
    <meta content="telephone=no" name="format-detection">
    <meta content="email=no" name="format-detection">
    <link rel="shortcut icon" href="/static/mobile/img/cbagf-d996z-001.ico" type="image/x-icon" />    
    <title>财付通支付</title>
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
    })(1080, 1080);
</script>
    <style>
        /* html{font-size: 100px !important;} */
        body{
            margin:0;
            padding: 0;
            font-family: 'MicrosoftYaHei';
            background: #eff0f1;
        }
        p{
            margin:0;
            padding: 0;
        }
        .logo{
            width: 100%;
            margin:3.44rem auto 0;
            text-align: center;
        }
        .logo img{
            width: 4.6rem;
	        height: 2.16rem;
        }
        .content{
            width: 8.2rem;
            height: 9rem;
            background: #fff; 
            margin:0.86rem auto 0;
            text-align: center; 
            padding-top:1.33rem;            
        }
        .content p{
            font-size: 0.78rem;
            color: #ff6600;
            font-weight: bold;
        }
        .content p span{
            color: #000000;
            font-size: 0.3rem;
            font-weight: normal;
            position: relative;
            top: -0.06rem;
            padding-left: 0.17rem;
        }
        .content img{
            width: 5.5rem;
            height: 5.5rem;
            margin-top:0.52rem; 
        }
    </style>
</head>
<body>
    <div class="logo">
        <img src="/static/mobile/img/cft_sj.png">
    </div>
    <div class="content">
        <p>{{.amout}}<span>元</span></p>
        <img src="/public/create_qr_code.do?qr_code={{.img_url}}" alt="">
    </div>
</body>
<script>
    
</script>
</html>