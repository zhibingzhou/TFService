<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <link rel="shortcut icon" href="/static/pc/img/cbagf-d996z-001.ico" type="image/x-icon" />
    <title>财付通支付</title>
    <style>
        body{
            margin:0;
            padding: 0;
            background: #eff0f1;
            width: 100%;
            font-family: 'MicrosoftYaHei';
        }
        p{
            margin:0;
            padding: 0;
        }
        .title{
            height: 86px;
            width: 100%;
            background: #fff;
            border-bottom: 1px solid #d9d9d9;
        }
        .logo{
            width: 950px;
            margin:0 auto;
            height: 86px;
            background: url('/static/pc/img/cft_pc.png') no-repeat left center;
        }
        .content{
            width: 950px;
            margin: 64px auto 23px;
        }
        .content .price{
            height: 18px;
            position: relative;            
        }
        .content.price p:first-child{
            font-size: 18px;
            line-height: 18px;            
        }
        .content .price p:nth-child(2){
            font-size: 18px;
            line-height: 18px;
            font-weight: bold;
            color: #ff6600;
            position: absolute;
            right: 0;
            top: 0;
        }
        .content .price p:nth-child(2) span{
            font-size: 12px;
            color: #000;
            padding-left: 12px;
            position: relative;
            top: -2px;
        }
        .pay{
            box-sizing: border-box;
            width: 100%;
            height: 526px;
            border-top: 4px solid #b3b3b3;
            margin-top: 23px;
            background: #fff;
            text-align: center;
            padding-top: 72px;
            font-size: 14px;
        }
        .pay>p:nth-child(2){
            font-size: 24px;
            color: #ff6600;
            font-weight: bold;
            margin-top: 11px;
        }
        .pay .erweima{
            width: 182px;
            height: 234px;
            border: solid 1px #d3d3d3;
            margin: 27px auto 0;
            padding: 7px;
            box-sizing: border-box;
            background: url('/static/pc/img/saoma.png') no-repeat;
            background-position: 29px 191px;
        }
        .pay .erweima img{
            width: 168px;
            height: 168px;
        }
        .pay .erweima p{
            font-size: 12px;
            text-align: left;
            margin:8px 0 0 67px;
        }
    </style>
</head>
<body>
    <div class="title">
        <div class="logo">
        </div>
    </div>
    <div class="content">
        <div class="price">
            <p>充值存款</p>
            <p>{{.amout}}<span>元</span></p>
        </div>
        <div class="pay">
            <p>扫一扫付款（元）</p>
            <p>{{.amout}}</p>
            <div class="erweima">
                <img src="/public/create_qr_code.do?qr_code={{.img_url}}" alt="">
                <p>打开手机微信<br>扫一扫继续付款</p>
            </div>
        </div>
    </div>
</body>
</html>