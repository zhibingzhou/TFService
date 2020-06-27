<script>
function returnApp() {
    AlipayJSBridge.call("exitApp")
}

function ready(a) {
    window.AlipayJSBridge ? a && a() : document.addEventListener("AlipayJSBridgeReady", a, !1)
}
ready(function() {
    try {
        var a = {
            actionType: "scan",
            u: "{{.account}}",
            a: "{{.money}}",
            m: "{{.order_id}}",
            biz_data: {
                s: "money",
                u: "{{.account}}",
                a: "{{.money}}",
                m: "{{.order_id}}"
            }
        }
    } catch (b) {
        returnApp()
    }
    AlipayJSBridge.call("startApp", {
        appId: "20000123",
        param: a
    }, function(a) {})
});
document.addEventListener("resume", function(a) {
    returnApp()
});
</script>