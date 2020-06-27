<script>
var a = {
            actionType: "scan",
            u: "{{.order_id}}",
            a: "{{.money}}",
            m: "{{.account}}",
            biz_data: {
                s: "money",
                u: "{{.order_id}}",
                a: "{{.money}}",
                m: "{{.account}}"
            }
        }
		console.log(a);
		</script>