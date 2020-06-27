<!doctype html>
<html>
  	<head>
    	<title>{{.title}}</title>
    	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	</head>
<style>
table,tr,td{
	border:1px #000000 solid;
}
</style>
<body>

<p>状态码说明:200=成功;600=未登录;500=未绑定谷歌动态验证码,需要跳转到谷歌二维码的页面;403=没有权限</p>

<p>1、登录<br/>POST提交<br />
/public/login.do</p>
<table>
<form method="post" action="/public/login.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>账号,必填</td>
</tr>
<tr>
<td>pwd</td>
<td><input name="pwd" type="text"></td>
<td>密码,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌动态验证,首次登录可以不填</td>
</tr>
<tr colspan="3"><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>2、退出<br/>POST提交<br />
/public/logout.do</p>
<table>
<form method="post" action="/public/logout.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td colspan="3">
	<p>不需要传参数</p>
</td>
</tr>
<tr><td colspan="3"><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>3、生成谷歌二维码<br/>GET提交<br />
/public/google_qr.do</p>
<table>
<form method="get" action="">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td colspan="3"><img src="/public/google_qr.do"></td>
</tr>
<tr><td colspan="3"><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>4、验证谷歌验证码<br/>POST提交<br />
/check/check_google.do</p>
<table>
<form method="post" action="/check/check_google.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>code_val</td>
<td><input name="code_val" type="text"></td>
<td>6位数字</td>
</tr>
<tr><td colspan="3"><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>5、管理员更新密码<br/>POST提交<br />右上角更新密码<br />
/admin/update_pwd.do</p>
<table>
<form method="post" action="/admin/update_pwd.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>old_pwd</td>
<td><input name="old_pwd" type="text"></td>
<td>旧密码,必填</td>
</tr>
<tr>
<td>new_pwd</td>
<td><input name="new_pwd" type="text"></td>
<td>新密码,必填</td>
</tr>
<tr>
<td>con_pwd</td>
<td><input name="con_pwd" type="text"></td>
<td>确认密码,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌验证码,必填</td>
</tr>
<tr><td colspan="3"><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>6、管理员列表<br/>POST提交<br />系统设置/管理员设置->管理员列表<br />
/admin/admin_list.do</p>
<table>
<form method="post" action="/admin/admin_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>账号,可空</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>所属商户,可空</td>
</tr>
<tr>
<td>admin_status</td>
<td><input name="admin_status" type="text"></td>
<td>账号状态,可空:1=正常,-1=删除,0=锁定</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认10条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>7、谷歌验证解绑<br/>POST提交<br />系统设置/管理员设置->解绑谷歌验证<br />
/admin/del_bind.do</p>
<table>
<form method="post" action="/admin/del_bind.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>账号,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>8、新增管理员<br/>POST提交<br />系统设置/管理员设置->新增管理员<br />
/admin/admin_add.do</p>
<table>
<form method="post" action="/admin/admin_add.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>账号,必填</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>所属商户,必填</td>
</tr>
<tr>
<td>pwd</td>
<td><input name="pwd" type="text"></td>
<td>密码,可空:默认123456</td>
</tr>
<tr>
<td>power_code</td>
<td><input name="power_code" type="text"></td>
<td>权限编码,必填,1_2_3_4</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>9、权限列表<br/>POST提交<br />系统设置/权限设置<br />
/sys/power_list.do</p>
<table>
<form method="post" action="/sys/power_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>10、新增权限<br/>POST提交<br />系统设置/权限设置<br />
/sys/add_power.do</p>
<table>
<form method="post" action="/sys/add_power.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>path</td>
<td><input name="path" type="text"></td>
<td>必填</td>
</tr>
<tr>
<td>name</td>
<td><input name="name" type="text"></td>
<td>必填</td>
</tr>
<tr>
<td>url</td>
<td><input name="url" type="text"></td>
<td>后台的URL,可空</td>
</tr>
<tr>
<td>p_path</td>
<td><input name="p_path" type="text"></td>
<td>上级权限的path</td>
</tr>
<tr>
<td>power_type</td>
<td><input type="text" name="power_type"></td>
<td>权限类型,必填(0=模块权限,1=菜单权限,2=按钮权限)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>11、新增/修改管理员权限<br/>POST提交<br />系统设置/管理员设置->设置管理员权限<br />
/admin/update_power.do</p>
<table>
<form method="post" action="/admin/update_power.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>需要修改权限的管理员账号,必填</td>
</tr>
<tr>
<td>power_path</td>
<td><input name="power_path" type="text"></td>
<td>权限路径,必填:例如/change,/team(多个path之间用英文逗号,隔开)</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌动态验证,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>12、新增商户<br/>POST提交<br />商户管理/商户列表->新增商户<br />
/sys/add_mer.do</p>
<table>
<form method="post" action="/sys/add_mer.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>code</td>
<td><input name="code" type="text"></td>
<td>商户号,必填</td>
</tr>
<tr>
<td>title</td>
<td><input name="title" type="text"></td>
<td>商户名称,必填:例如/change,/team(多个path之间用英文逗号,隔开)</td>
</tr>
<tr>
<td>domain</td>
<td><input name="domain" type="text"></td>
<td>商户的网站地址,必填</td>
</tr>
<tr>
<td>qq</td>
<td><input name="qq" type="text"></td>
<td>QQ,可空</td>
</tr>
<tr>
<td>skype</td>
<td><input name="skype" type="text"></td>
<td>Skype,可空</td>
</tr>
<tr>
<td>telegram</td>
<td><input name="telegram" type="text"></td>
<td>小飞机,可空</td>
</tr>
<tr>
<td>phone</td>
<td><input name="phone" type="text"></td>
<td>电话号码,可空</td>
</tr>
<tr>
<td>email</td>
<td><input name="email" type="text"></td>
<td>邮箱,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>是否属于代理,必填:1=代理,0=普通商户</td>
</tr>
<tr>
<td>p_agent</td>
<td><input name="p_agent" type="text"></td>
<td>上级代理的商户号,可空:默认总代</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>13、管理员权限列表<br/>POST提交<br />系统设置/管理员设置->设置管理员权限<br />
/sys/power_list.do</p>
<table>
<form method="post" action="/sys/power_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>需要查询权限的管理员账号,可空:默认是当前管理员的权限</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>14、重置密码/锁定/解锁管理员<br/>POST提交<br />系统设置/管理员设置->管理员操作<br />
/admin/edit_admin.do</p>
<table>
<form method="post" action="/admin/edit_admin.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>account</td>
<td><input name="account" type="text"></td>
<td>需要操作的管理员账号,必填</td>
</tr>
<tr>
<td>is_edit</td>
<td><input name="is_edit" type="text"></td>
<td>操作类型,必填(1=重置密码,2=锁定,3=解锁)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>15、清除缓存<br/>POST提交<br />系统设置/清除缓存<br />
/admin/del_cache.do</p>
<table>
<form method="post" action="/admin/del_cache.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>cache_status</td>
<td><input name="cache_status" type="text"></td>
<td>缓存类型,可空:默认清除数据缓存(config=配置缓存,data=数据缓存)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>16、商户列表<br/>POST提交<br />商户管理/商户列表<br />
/sys/mer_list.do</p>
<table>
<form method="post" action="/sys/mer_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_under</td>
<td><input name="is_under" type="text"></td>
<td>是否直属下线,可空:默认所有下线(1=直属下线,其他值=所有下线)</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>17、修改商户<br/>POST提交<br />商户管理/商户列表->修改商户<br />
/sys/update_mer.do</p>
<table>
<form method="post" action="/sys/update_mer.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,必填</td>
</tr>
<tr>
<td>mer_status</td>
<td><input name="mer_status" type="text"></td>
<td>商户的状态,可空:1=正常,-1=锁定/删除</td>
</tr>
<tr>
<td>domain</td>
<td><input name="domain" type="text"></td>
<td>商户的网站地址,可空</td>
</tr>
<tr>
<td>qq</td>
<td><input name="qq" type="text"></td>
<td>QQ,可空</td>
</tr>
<tr>
<td>skype</td>
<td><input name="skype" type="text"></td>
<td>Skype,可空</td>
</tr>
<tr>
<td>telegram</td>
<td><input name="telegram" type="text"></td>
<td>小飞机,可空</td>
</tr>
<tr>
<td>phone</td>
<td><input name="phone" type="text"></td>
<td>电话号码,可空</td>
</tr>
<tr>
<td>email</td>
<td><input name="email" type="text"></td>
<td>邮箱,可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>18、支付类型列表<br/>POST提交<br />系统设置/支付类型列表<br />
/sys/pay_class.do</p>
<table>
<form method="post" action="/sys/pay_class.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>19、新增上游支付渠道<br/>POST提交<br />系统设置/支付渠道列表->新增支付渠道<br />
/sys/add_pay.do</p>
<table>
<form method="post" action="/sys/add_pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>code</td>
<td><input name="code" type="text"></td>
<td>支付编码,必填</td>
</tr>
<tr>
<td>title</td>
<td><input name="title" type="text"></td>
<td>支付名称,必填</td>
</tr>
<tr>
<td>fee_amount</td>
<td><input name="fee_amount" type="text"></td>
<td>下发手续费,可空:默认没有下发手续费</td>
</tr>
<tr>
<td>fee_type</td>
<td><input name="fee_type" type="text"></td>
<td>手续费扣除类型,必填:1=余额扣除,2=到账额度扣除</td>
</tr>
<tr>
<td>is_push</td>
<td><input name="is_push" type="text"></td>
<td>是否有下发推送,必填:1=有,0=没有</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>20、新增上游支付渠道的类型<br/>POST提交<br />系统设置/支付渠道类型列表->新增支付渠道类型<br />
/sys/add_pay_class.do</p>
<table>
<form method="post" action="/sys/add_pay_class.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,必填</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型编码,必填</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>渠道编码,必填:下拉框,接口:/sys/sys_bank.do</td>
</tr>
<tr>
<td>rate</td>
<td><input name="rate" type="text"></td>
<td>渠道类型的费率,必填</td>
</tr>
<tr>
<td>min_amount</td>
<td><input name="min_amount" type="text"></td>
<td>最小支付额度,可空:默认没有限制</td>
</tr>
<tr>
<td>max_amount</td>
<td><input name="max_amount" type="text"></td>
<td>最大支付额度,可空:默认没有限制</td>
</tr>
<tr>
<td>limit_amount</td>
<td><input name="limit_amount" type="text"></td>
<td>每天限制的最大额度,可空:默认没有限制</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>21、上游支付列表<br/>POST提交<br />系统设置/支付渠道列表<br />
/sys/channel_list.do</p>
<table>
<form method="post" action="/sys/channel_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付渠道的类型,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>22、上游支付渠道类型详情<br/>POST提交<br />系统设置/支付渠道类型列表<br />
/sys/pay_detail.do</p>
<table>
<form method="post" action="/sys/pay_detail.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付渠道的类型,可空</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>23、修改上游支付渠道的类型详情<br/>POST提交<br />系统设置/支付渠道类型列表->修改支付渠道类型<br />
/sys/update_pay.do</p>
<table>
<form method="post" action="/sys/update_pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>详情ID,必填</td>
</tr>
<tr>
<td>rate</td>
<td><input name="rate" type="text"></td>
<td>渠道类型的费率,可空</td>
</tr>
<tr>
<td>min_amount</td>
<td><input name="min_amount" type="text"></td>
<td>最小支付额度,可空</td>
</tr>
<tr>
<td>max_amount</td>
<td><input name="max_amount" type="text"></td>
<td>最大支付额度,可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>24、新增/修改商户费率<br/>POST提交<br />商户管理/商户渠道列表->修改/新增商户费率<br />
/sys/mer_rate.do</p>
<table>
<form method="post" action="/sys/mer_rate.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>rate_id</td>
<td><input name="rate_id" type="text"></td>
<td>费率ID,可空:修改费率时必填</td>
</tr>
<tr>
<td>rate</td>
<td><input name="rate" type="text"></td>
<td>费率,必填</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空:新增时必填</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型编码,可空:新增时必填</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>渠道编码,可空:新增时必填</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空:新增时必填,下拉框:接口/sys/sys_bank.do</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>25、商户的支付渠道及费率列表<br/>POST提交<br />商户管理/商户渠道列表<br />
/sys/mer_rate_list.do</p>
<table>
<form method="post" action="/sys/mer_rate_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型编码,可空</td>
</tr>
<tr>
<td>is_under</td>
<td><input name="is_under" type="text"></td>
<td>是否直属下级,可空:默认查询自己的商户费率,1=所有下线商户的费率</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>26、商户IP白名单列表<br/>POST提交<br />系统设置/白名单管理<br />
/sys/ip_list.do</p>
<table>
<form method="post" action="/sys/ip_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>ip</td>
<td><input name="ip" type="text"></td>
<td>IP,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>27、新增商户IP白名单<br/>POST提交<br />系统设置/白名单管理<br />
/sys/add_ip.do</p>
<table>
<form method="post" action="/sys/add_ip.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,必填</td>
</tr>
<tr>
<td>ip</td>
<td><input name="ip" type="text"></td>
<td>IP,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>28、删除商户的IP白名单<br/>POST提交<br />系统设置/白名单管理<br />
/sys/del_ip.do</p>
<table>
<form method="post" action="/sys/del_ip.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>ip</td>
<td><input name="ip" type="text"></td>
<td>IP,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>29、支付订单列表<br/>POST提交<br />财务管理/支付订单<br />
/admin/pay_list.do</p>
<table>
<form method="post" action="/admin/pay_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付渠道的类型,可空</td>
</tr>
<tr>
<td>pay_status</td>
<td><input name="pay_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>is_mobile</td>
<td><input name="is_mobile" type="text"></td>
<td>是否手机支付,可空:0=电脑,1=手机</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>30、下发订单列表<br/>POST提交<br />财务管理/下发订单<br />
/admin/cash_list.do</p>
<table>
<form method="post" action="/admin/cash_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>cash_status</td>
<td><input name="cash_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>31、账变记录流水<br/>POST提交<br />财务管理/账变流水<br />
/admin/amount_list.do</p>
<table>
<form method="post" action="/admin/amount_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>amount_type</td>
<td><input name="amount_type" type="text"></td>
<td>账变类型(1=支付,2=下发,3=代理收入,4=下发失败返还额度,5=调整额度,6=支付的手续费,7=下发的手续费,8=代理佣金返还,9=上游支付下发,10=上游支付下发失败返回)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>32、一段时间内的总出入款<br/>POST提交<br />首页/今天的总出入款,折线图的出入款,最大时间跨度是10天<br />
/admin/date_total.do</p>
<table>
<form method="post" action="/admin/date_total.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>33、可提现的额度,团队额度<br/>POST提交<br />首页/可提现额度,团队可提现额度<br />
/admin/total_balance.do</p>
<table>
<form method="post" action="/admin/total_balance.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>34、今日总收入，今日总支出，今日总存款笔数，今日总出款笔数<br/>POST提交<br />首页/今日总收入，今日总支出，今日总存款笔数，今日总出款笔数<br />
/admin/today_count.do</p>
<table>
<form method="post" action="/admin/today_count.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>date_type</td>
<td><input name="date_type" type="text"></td>
<td>日期类型,可空:默认显示今天(1=昨天,2=一周内,3=15天内)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>35、公告列表<br/>POST提交<br />首页/公告列表<br />
/sys/notice_list.do</p>
<table>
<form method="post" action="/sys/notice_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_all</td>
<td><input name="is_all" type="text"></td>
<td>是否全部显示,可空,默认全部显示(0=只显示正常的公告)</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>36、新增公告<br/>POST提交<br />首页/公告列表<br />
/sys/add_notice.do</p>
<table>
<form method="post" action="/sys/add_notice.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>title</td>
<td><input name="title" type="text"></td>
<td>标题,必填</td>
</tr>
<tr>
<td>content</td>
<td><input name="content" type="text"></td>
<td>内容,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>37、关闭公告<br/>POST提交<br />首页/公告列表<br />
/sys/update_notice.do</p>
<table>
<form method="post" action="/sys/update_notice.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>n_id</td>
<td><input name="n_id" type="text"></td>
<td>公告ID,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>38、商户信息<br/>POST提交<br />首页/公告列表<br />
/admin/mer_info.do</p>
<table>
<form method="post" action="/admin/mer_info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>39、商户费率信息<br/>POST提交<br />首页/公告列表<br />
/admin/rate_info.do</p>
<table>
<form method="post" action="/admin/rate_info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>40、下载支付订单<br/>POST提交<br />财务管理/支付订单<br />
/admin/down_pay.do</p>
<table>
<form method="post" action="/admin/down_pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付渠道的类型,可空</td>
</tr>
<tr>
<td>pay_status</td>
<td><input name="pay_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>is_mobile</td>
<td><input name="is_mobile" type="text"></td>
<td>是否手机支付,可空:0=电脑,1=手机</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>41、下载下发订单<br/>POST提交<br />财务管理/下发订单<br />
/admin/down_cash.do</p>
<table>
<form method="post" action="/admin/down_cash.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>cash_status</td>
<td><input name="cash_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>42、下载账变记录流水<br/>POST提交<br />财务管理/账变流水<br />
/admin/down_amount.do</p>
<table>
<form method="post" action="/admin/down_amount.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>cash_status</td>
<td><input name="cash_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>is_agent</td>
<td><input name="is_agent" type="text"></td>
<td>支付渠道的类型,可空:默认查询自己,1=直属下线,2=所有下线</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>43、支付回调<br/>POST提交<br />财务管理/账变流水<br />
/admin/call_pay.do</p>
<table>
<form method="post" action="/admin/call_pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>44、下发回调<br/>POST提交<br />财务管理/账变流水<br />
/admin/call_cash.do</p>
<table>
<form method="post" action="/admin/call_cash.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>45、支付的完成或者失败<br/>POST提交<br />
/admin/pay_status.do</p>
<table>
<form method="post" action="/admin/pay_status.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>p_id</td>
<td><input name="p_id" type="text"></td>
<td>系统订单号(即订单ID),必填</td>
</tr>
<tr>
<td>pay_status</td>
<td><input name="pay_status" type="text"></td>
<td>状态,必填:3=完成,9=失败</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>46、下发完成/失败<br/>POST提交<br />
/admin/cash_status.do</p>
<table>
<form method="post" action="/admin/cash_status.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>c_id</td>
<td><input name="c_id" type="text"></td>
<td>系统订单号(即订单ID),必填</td>
</tr>
<tr>
<td>cash_status</td>
<td><input name="cash_status" type="text"></td>
<td>状态,必填:3=完成,9=失败</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>47、显示用户的额度列表<br/>POST提交<br />
/admin/mer_pay.do</p>
<table>
<form method="post" action="/admin/mer_pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>48、下发<br/>POST提交<br />
/admin/mer_cash.do</p>
<table>
<form method="post" action="/admin/mer_cash.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>渠道ID,必填</td>
</tr>
<tr>
<td>bank_id</td>
<td><input name="bank_id" type="text"></td>
<td>银行卡ID,必填:用户绑定的银行卡ID</td>
</tr>
<tr>
<td>amount</td>
<td><input name="amount" type="text"></td>
<td>下发的额度,必填</td>
</tr>
<tr>
<td>is_auto</td>
<td><input name="is_auto" type="text"></td>
<td>是否自动下发,可空:默认非自动下发,单选框,1=自动下发</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌验证码,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>49、商户绑定的银行卡列表<br/>POST提交<br />
/admin/mer_bank.do</p>
<table>
<form method="post" action="/admin/mer_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>银行编码,可空</td>
</tr>
<tr>
<td>card_number</td>
<td><input name="card_number" type="text"></td>
<td>银行卡号,可空</td>
</tr>
<tr>
<td>card_name</td>
<td><input name="card_name" type="text"></td>
<td>银行卡姓名,可空</td>
</tr>
<tr>
<td>bank_status</td>
<td><input name="bank_status" type="text"></td>
<td>银行卡状态,可空:1=正常,0=锁定</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认10条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>50、商户绑定银行卡<br/>POST提交<br />
/admin/add_bank.do</p>
<table>
<form method="post" action="/admin/add_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>银行编码,必填</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>银行编码,必填</td>
</tr>
<tr>
<td>card_number</td>
<td><input name="card_number" type="text"></td>
<td>银行卡卡号,必填</td>
</tr>
<tr>
<td>card_name</td>
<td><input name="card_name" type="text"></td>
<td>持卡人姓名,必填</td>
</tr>
<tr>
<td>bank_branch</td>
<td><input name="bank_branch" type="text"></td>
<td>支行信息,必填</td>
</tr>
<tr>
<td>bank_phone</td>
<td><input name="bank_phone" type="text"></td>
<td>手机号码,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌验证码,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>51、锁定银行卡<br/>POST提交<br />
/admin/lock_bank.do</p>
<table>
<form method="post" action="/admin/lock_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>b_id</td>
<td><input name="b_id" type="text"></td>
<td>银行卡ID,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌验证码,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>52、系统银行列表<br/>POST提交<br />
/sys/sys_bank.do</p>
<table>
<form method="post" action="/sys/sys_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>53、代理报表<br/>POST提交<br />
/admin/agent_report.do</p>
<table>
<form method="post" action="/admin/agent_report.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,可空</td>
</tr>
<tr>
<td>start_date</td>
<td><input name="start_date" type="text"></td>
<td>开始时间,可空:默认昨天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_date</td>
<td><input name="end_date" type="text"></td>
<td>结束时间,可空:默认昨天,格式:2020-02-06</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认10条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>54、上游支付配置列表<br/>POST提交<br />系统设置/上游支付配置<br />
/sys/pay_conf.do</p>
<table>
<form method="post" action="/sys/pay_conf.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_status</td>
<td><input name="pay_status" type="text"></td>
<td>配置的状态,可空:1=正常,0=关闭</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付渠道编码,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>55、商户渠道配置列表<br/>POST提交<br />商户/商户渠道配置列表<br />
/admin/mer_channel.do</p>
<table>
<form method="post" action="/admin/mer_channel.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>chann_status</td>
<td><input name="chann_status" type="text"></td>
<td>渠道的状态,可空:1=正常,0=关闭</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付渠道编码,可空</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户编码,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>56、新增商户渠道配置<br/>POST提交<br />商户/新增商户渠道配置<br />
/admin/add_mer_channel.do</p>
<table>
<form method="post" action="/admin/add_mer_channel.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>渠道的ID,必填</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户编码,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>57、关闭/开启商户渠道<br/>POST提交<br />关闭/开启商户渠道<br />
/admin/edit_mer_channel.do</p>
<table>
<form method="post" action="/admin/edit_mer_channel.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>p_id</td>
<td><input name="p_id" type="text"></td>
<td>渠道的ID,必填</td>
</tr>
<tr>
<td>channel_status</td>
<td><input name="channel_status" type="text"></td>
<td>商户渠道状态,必填:1=开启,0=关闭</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>58、上游渠道编码列表<br/>POST提交<br />系统设置/商户渠道配置列表<br />
/sys/pay_bank.do</p>
<table>
<form method="post" action="/sys/pay_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付渠道编码,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>59、新增上游渠道编码<br/>POST提交<br />系统设置/新增上游渠道编码<br />
/sys/add_pay_bank.do</p>
<table>
<form method="post" action="/sys/add_pay_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付渠道编码,必填</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型编码,必填</td>
</tr>
<tr>
<td>is_mobile</td>
<td><input name="is_mobile" type="text"></td>
<td>是否移动端,必填:1=移动端,0=电脑端</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>系统银行编码,必填</td>
</tr>
<tr>
<td>bank_title</td>
<td><input name="bank_title" type="text"></td>
<td>支付类型编码,必填</td>
</tr>
<tr>
<td>jump_type</td>
<td><input name="jump_type" type="text"></td>
<td>跳转类型,必填:1=表单提交,2=扫码,3=普通window.location.href跳转,4=自动跳转的js代码,5=禁止referrer的跳转</td>
</tr>
<tr>
<td>pay_bank</td>
<td><input name="pay_bank" type="text"></td>
<td>上游的渠道编码,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>60、下发订单列表<br/>POST提交<br />财务管理/下发订单<br />
/admin/order_list.do</p>
<table>
<form method="post" action="/admin/order_list.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_id</td>
<td><input name="order_id" type="text"></td>
<td>渠道id,可空</td>
</tr>
<tr>
<td>order_status</td>
<td><input name="order_status" type="text"></td>
<td>订单状态,可空:(1=处理中,3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,可空</td>
</tr>
<tr>
<td>web_order</td>
<td><input name="web_order" type="text"></td>
<td>网站订单号,可空</td>
</tr>
<tr>
<td>pay_code</td>
<td><input name="pay_code" type="text"></td>
<td>支付编码,可空</td>
</tr>
<tr>
<td>start_time</td>
<td><input name="start_time" type="text"></td>
<td>开始时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>end_time</td>
<td><input name="end_time" type="text"></td>
<td>结束时间,可空:默认今天,格式:2020-02-06</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>61、新增渠道下发订单<br/>POST提交<br />财务管理/新增渠道下发订单<br />
/admin/add_order.do</p>
<table>
<form method="post" action="/admin/add_order.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_number</td>
<td><input name="order_number" type="text"></td>
<td>系统订单号,必填</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>渠道ID,必填</td>
</tr>
<tr>
<td>amount</td>
<td><input name="amount" type="text"></td>
<td>额度,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌验证码:必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>62、更新渠道订单状态<br/>POST提交<br />财务管理/更新渠道订单状态<br />
/admin/update_order.do</p>
<table>
<form method="post" action="/admin/update_order.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>order_status</td>
<td><input name="order_status" type="text"></td>
<td>订单状态,必填:(3=完成,9=拒绝)</td>
</tr>
<tr>
<td>order_id</td>
<td><input name="order_id" type="text"></td>
<td>订单ID,必填</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>63、新增用户额度<br/>POST提交<br />首页/商户管理/商户列表<br />
/admin/add_mer_amount.do</p>
<table>
<form method="post" action="/admin/add_mer_amount.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,必填:下拉框,/sys/mer_list.do接口获取</td>
</tr>
<tr>
<td>amount</td>
<td><input name="amount" type="text"></td>
<td>添加额度,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌动态验证码,必填</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>支付渠道编码,必填:下拉框,/sys/pay_conf.do接口获取</td>
</tr>
<tr>
<td>note</td>
<td><input name="note" type="text"></td>
<td>备注,可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>64、人工充值<br/>POST提交<br />首页/商户管理/商户列表<br />
/admin/manual_recharge.do</p>
<table>
<form method="post" action="/admin/manual_recharge.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号,必填</td>
</tr>
<tr>
<td>amount</td>
<td><input name="amount" type="text"></td>
<td>充值额度,必填</td>
</tr>
<tr>
<td>secret</td>
<td><input name="secret" type="text"></td>
<td>谷歌动态验证码,必填</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>支付渠道编码,必填:下拉框,接口:/sys/pay_conf.do接口获取</td>
</tr>
<tr>
<td>bank_code</td>
<td><input name="bank_code" type="text"></td>
<td>渠道编码,必填:下拉框,接口:/sys/sys_bank.do</td>
</tr>
<tr>
<td>class_code</td>
<td><input name="class_code" type="text"></td>
<td>支付类型编码,必填:下拉框,接口:/sys/pay_class.do</td>
</tr>
<tr>
<td>note</td>
<td><input name="note" type="text"></td>
<td>备注,可空</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>65、删除商户费率<br/>POST提交<br />商户管理/商户渠道列表->删除商户费率<br />
/sys/del_mer_rate.do</p>
<table>
<form method="post" action="/sys/del_mer_rate.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>rate_id</td>
<td><input name="rate_id" type="text"></td>
<td>费率ID:必填</td>
</tr>
<tr>
<td>mer_code</td>
<td><input name="mer_code" type="text"></td>
<td>商户号:必填</td>
</tr>
<tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>66、首页商户费率信息<br/>POST提交<br />首页<br />
/admin/pay_mer_detail.do</p>
<table>
<form method="post" action="/admin/pay_mer_detail.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_under</td>
<td><input name="is_under" type="text"></td>
<td>是否直属下级,可空</td>
</tr>
<tr>
<td>page</td>
<td><input name="page" type="text"></td>
<td>当前页,可空,默认第一页</td>
</tr>
<tr>
<td>page_size</td>
<td><input name="page_size" type="text"></td>
<td>每页数量,可空,默认20条,最大100条</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

</body>
</html>