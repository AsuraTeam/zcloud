/**
 * Created by zhaoyun on 2017/12/31.
 */
/**
 *
 * @param paramter
 * @param url
 * @return {string}
 */
function post(paramter, url) {
    var result = "";
    $.ajax({
        type: "POST",
        url: url,
        data: paramter,
        async: false,
        success: function (data) {
            result = data;
        }
    });
    return result;
}
function del(paramter, url) {
    var result = "";
    $.ajax({
        type: "DELETE",
        url: url,
        data: paramter,
        async: false,
        success: function (data) {
            result = data;
        }
    });
    return result;
}
function get(paramter, url) {
    var result = "";
    $.ajax({
        type: "GET",
        url: url,
        data: paramter,
        async: false,
        success: function (data) {
            result = data;
        }
    });
    return result;
}
//获取from数据
function get_form_data() {
    var result = {};
    var forch = ["input", "textarea", "select"];
    for (i = 0; i < forch.length; i++) {
        $.each($("form " + forch[i]),
            function (name, object) {
                result[$(object).attr("name")] = $(object).val()
            }
        );
    }
    return result;
}

function loginOut() {
    var url = "/api/user/logout";
    post({}, url);
    setTimeout(function () {
        window.location.href = "/login"
    }, 2000)
}


function checkChange(keys, type) {
    if (!type) {
        type = ""
    }
    var data = get_form_data();
    console.log(data);
    checkValue(data, keys, type)
}

/**
 * 检查yaml语法是否正确
 * @param obj
 */
function checkYaml(data) {
    var url = "/api/template/yaml/check";
    var result = post({yaml: data["Yaml"]}, url);
    if (result == "true") {
        return true
    }
    return false
}


/**
 * 检查docker安装路径
 * @param str
 */
function checkDockerInstallDir(str) {
    if (str.length < 3) {
        return false
    }
    var regEx = /^([\/][\w-]+)*$/i;
    if (regEx.test(str)) {
        return true
    }
    return false
}


/**
 * 检查版本输入正确性
 * @param str
 * @return {boolean}
 */
function checkImageTag(str) {

    if (str.length < 3) {
        return false
    }
    var regEx =  /^[0-9a-zA-Z]*$/i;
    if (regEx.test(str)) {
        return true
    }
    return false
}


/**
 * 检查IP地址是否正确
 * @param str
 * @returns {boolean}
 */
function checkIp(str) {
    var exp = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
    var reg = str.match(exp);
    if (reg) {
        return true
    }
    return false
}
/**
 * 检查集群名称是否正常
 * @param str
 * @returns {boolean}
 */
function checkClusterName(str) {
    var regEx = /^[a-zA-Z_]\w{2,35}/;
    if (regEx.test(str)) {
        return true
    }
    return false
}

/**
 * 2018-01-22 10:14
 * 检查设置仓库权限时项目是否正确,不能已 / 开头和结尾
 * @param str
 * @return {boolean}
 */
function checkProject(str) {
    var regEx = /^\/(.*)\/$/;
    if(regEx.test(str)){
        return false
    }
    regEx = /^\/(.*)/;
    if(regEx.test(str)){
        return false
    }
    regEx = /(.*)\/$/;
    if(regEx.test(str)){
        return false
    }
    console.log(str)
    return true
}

/**
 * 检查集群名称是否正常
 * @param str
 * @returns {boolean}
 */
function checkAppName(str) {
    var regEx = /^[a-zA-Z]\w{2,35}/;
    if (regEx.test(str)) {
        return true
    }
    return false
}

/**
 * 检查应用名称和服务名称
 * @param str
 * @return {boolean}
 */
function checkAppNameService(str) {
    var serviceName = $('input[name="ServiceName"]').val();
    var cluster = $('input[name="ClusterName"]').val();
    var regEx = /^[a-zA-Z]\w{2,35}/;
    if (regEx.test(str)) {
        var url = "/api/service/name?AppName=" + str + "&ClusterName=" + cluster+"&ServiceName="+serviceName;
        var data = get({}, url);
        if (data.length == 0) {
            return true;
        }
        return false
    }

    return false
}

/**
 * 检查配额名称是否正确
 * @param str
 * @return {boolean}
 */
function checkQuotaName(str) {
    var regEx = /^[a-zA-Z_]\w{4,32}/;
    if (regEx.test(str)) {
        return true
    }
    return false
}

/**
 * 检查模板名称是否正确
 * @param str
 * @return {boolean}
 */
function checkTemplateName(str) {
    var regEx = /^[a-zA-Z_]\w{4,32}/;
    if (regEx.test(str)) {
        return true
    }
    return false
}

/**
 * 检查集群网卡名称
 * @param str
 * @returns {boolean}
 */
function checkNetworkCart(str) {
    var regEx = /^[a-zA-Z]\w{2,10}/;
    if (regEx.test(str)) {
        return true
    }
    return false
}


function setInputOk(obj) {
    obj.css("font-size", "14px");
    obj.css("border", "1px solid #f0f0f0");
    obj.css("background-color", "#ffffff");
    obj.removeAttr("title")
}

function setInputError(obj, key) {
    if (!obj) {
        return;
    }
    obj.css("border", "1px solid #f16a7c");
    obj.css("background-color", "#fffff6");
    obj.css("font-size", "12px");

    var err = obj.attr(key);
    if (!err && key == "nullmsg") {
        err = "该项目必须填写"
    }
    obj.attr("placeholder", err);
    obj.attr("title", err)
}

/**
 * 验证数据中内容是否符合规定
 * 需要在验证标签 有 errmsg nullmsg validFunc 3个参数
 * errmsg 执行validFunc验证失败时提示信息
 * nullmsg 没有写入时错误信息
 * validFunc 要执行的自定义检查函数
 * @param data
 * @param keys
 * @returns {boolean}
 */
function checkValue(data, keys, type) {
    if (!type) {
        type = "input"
    }
    keys = keys.split(",");
    for (var i = 0; i < keys.length; i++) {
        var b = $(type + "[name='" + keys[i] + "']");

        if (data[keys[i]]) {
            var checkFunc = b.attr("validFunc");
            if (checkFunc) {
                var evalfun = checkFunc + "('" + data[keys[i]] + "')";

                var check = eval(evalfun);
                if (check) {
                    setInputOk(b)
                } else {
                    console.log(b);
                    setInputError(b, "errmsg");
                    return false
                }
            } else {
                setInputOk(b)
            }
            continue
        } else {
            console.log(b);
            setInputError(b, "nullmsg");
            return false
        }
    }
    return true
}

function success(msg) {
    $("#success-info-html-id").show();
    $("#success-info-id").html(msg);
    setTimeout(function () {
        $("#success-info-html-id").hide();
    }, 5000)
}

function faild(msg) {
    $("#faild-info-html-id").show();
    $("#faild-info-id").html(msg);
    setTimeout(function () {
        $("#faild-info-html-id").hide();
    }, 5000)
}

/**
 * Theme: Adminox Template
 * Author: Coderthemes
 * SweetAlert
 */

!function ($) {
    "use strict";

    var SweetAlert = function () {
    };

    //examples
    SweetAlert.prototype.init = function () {
        //Parameter
        $('#login-out').click(function () {
            swal({
                title: '退出登录',
                text: "",
                type: 'warning',
                showCancelButton: true,
                confirmButtonText: '马上退出',
                cancelButtonText: '再待一会',
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                swal(
                    '退出登录!',
                    '您已退出登录.',
                    'success'
                )
                loginOut()
            }, function (dismiss) {
                // dismiss can be 'cancel', 'overlay',
                // 'close', and 'timer'
            })
        });
    },
        //init
        $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
}(window.jQuery),

//initializing
    function ($) {
        "use strict";
        $.SweetAlert.init()
    }(window.jQuery);


jQuery.cookie = function (name, value, options) {
    if (typeof value != 'undefined') {
        options = options || {};
        if (value === null) {
            value = '';
            options = $.extend({}, options);
            options.expires = -1;
        }
        var expires = '';
        if (options.expires && (typeof options.expires == 'number' || options.expires.toUTCString)) {
            var date;
            if (typeof options.expires == 'number') {
                date = new Date();
                date.setTime(date.getTime() + (options.expires * 24 * 60 * 60 * 1000));
            } else {
                date = options.expires;
            }
            expires = '; expires=' + date.toUTCString();
        }
        options
        var path = options.path ? '; path=' + (options.path) : '';
        var path = ";path=/";
        var domain = options.domain ? '; domain=' + (options.domain) : '';
        var secure = options.secure ? '; secure' : '';
        document.cookie = [name, '=', encodeURIComponent(value), expires, path, domain, secure].join('');
    } else {
        var cookieValue = null;
        if (document.cookie && document.cookie != '') {
            var cookies = document.cookie.split(';');
            for (var i = 0; i < cookies.length; i++) {
                var cookie = jQuery.trim(cookies[i]);
                if (cookie.substring(0, name.length + 1) == (name + '=')) {
                    cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
                    break;
                }
            }
        }
        return cookieValue;
    }
};

function setBorderMove(obj) {
    obj.css("border", "1px solid #24a7e3")
}

function setBorderOut(obj) {
    obj.css("border", "1px solid #e0e0e0")
}

/**
 * 获取cookie中的cluster名称
 * @return
 */
function getClusterName(alias) {
    if (alias) {
        return $.cookie("clusterAlias")
    }
    return $.cookie("cluster")
}


/**
 * 加载集群在页面的集群选择信息
 * 2018-01-04
 */
function loadClusterSelect(id,nil, select) {
    var url = "/api/cluster/name";
    var html = ""
    if (nil){
        html += "<option>--请选择--</option>"
    }
    if(select){
        html += "<option value='"+select+"'>"+getClusterAlias(select)+"</option>";
    }
    var result = get({}, url);
    //  先获取cookie选择好的
    var cluster = getClusterName();
    for (var i = 0; i < result.length; i++) {
        if (!cluster) {
            cluster = setSelectCluster(result[i]["ClusterName"])
        }
        html += "<option value='" + result[i]["ClusterName"] + "'>" + result[i]["ClusterAlias"] + "</option>"
    }
    if (html.length < 5) {
        html = "<option>还没有集群哦</option>"
    }
    $("#"+id).html(html);
    $("#cluster_data").val(JSON.stringify(result))
}


/**
 * 获取集群的别名
 * 2018-01-22 11:44
 * @param clusterName
 * @return {string}
 */
function getClusterAlias(clusterName) {
    var data =  $("#cluster_data").val();
    var data = JSON.parse(data);
    for(var i=0;i<data.length;i++){
       if( data[i]["ClusterName"] == clusterName ){
           return data[i]["ClusterAlias"];
       }
    }
    return ""
}

/**
 * 设置选择集群后的cookie信息
 * @return
 */
function setSelectCluster(clusterName, clusterAlias) {
    $.cookie("cluster", clusterName)
    $.cookie("clusterAlias", clusterAlias)
    return clusterName
}

/**
 * 检查端口是否正确,在添加服务时使用
 * @param val
 * @return {boolean}
 */
function checkPort(data, name) {
    console.log(data)
    if (!data[name]) {
        return false
    }
    var vals = data[name].split(",")
    console.log(vals)
    for (var i = 0; i < vals.length; i++) {
        if (isNaN(vals[i])) {
            return false;
        }
        if (vals[i] > 65535 || vals[i] < 1) {
            return false
        }
    }
    return true
}

/**
 * 检查是否是数据
 * @param data
 * @param name
 * @return {boolean}
 */
function checkNumber(data, name) {
    if (isNaN(data[name])) {
        console.log(11111)
        return false;
    }
    return true
}

/**
 * 查询资源配额名称
 * @param quotaName
 * @param htmlId
 * @param isLock
 */
function setQuotaName(quotaId, htmlId, isLock) {
    var url = "/api/quota/name";
    console.log(url);
    var data = get({}, url);
    console.log(data)
    var html = "";
    for (var i = 0; i < data.length; i++) {
        if (data[i]["QuotaId"] + "" == quotaId + "" || data[i]["QuotaName"] == quotaId ) {
            html = "<option value='" + data[i]["QuotaName"] + "'>" + data[i]["QuotaName"] + "</option>";
            break
        }
    }
    console.log(html);

    if(isLock){
        $("#" + htmlId).html(html);
        return;
    }
    console.log(html);
    for (var i = 0; i < data.length; i++) {
        if (data[i]["QuotaName"] + "" != quotaId + ""  ) {
            html += "<option value='" + data[i]["QuotaName"] + "'>" + data[i]["QuotaName"] + "</option>";
        }
    }
    $("#" + htmlId).html(html);
}


/**
 * 设置模板数据
 * 2018-01-05
 * @param htmlId
 */
function setTemplateData(htmlId) {
    var url = "/api/template/name"
    var data = get({}, url)
    var html = "<option value=''>--</option>"
    for (var i = 0; i < data.length; i++) {
        html += "<option value='" + data[i]["TemplateId"] + "'>" + data[i]["TemplateName"] + "</option>"
    }
    $("#" + htmlId).html(html)
}


// 设置box高低
function setBoxHeight(id) {
    var active = $("#" + id).attr("class");
    if (active.indexOf("active") != -1) {
        return
    }
    var height = $("#" + id).outerHeight()
    height = height + 400;
    setTimeout(function () {
        $("#cart-box").css("height", height + "px")
    }, 100)
}


$(".all-select").change(function () {
    var check = $(this).is(":checked")
    $(".all").prop("checked", check)
    if (check) {
        $(".button-select").removeAttr("disabled")
    } else {
        $(".button-select").prop("disabled", "true")
    }
});

/**
 * checkbox 修改状态后操作
 * @param obj
 */
function checkBoxChange(obj) {
    var check = $(obj).is(":checked")
    if (check) {
        $(".button-select").removeClass("all-select")
        $(".button-select").removeAttr("disabled")
    } else {
        $(".button-select").prop("disabled", "true")
        $(".button-select").addClass("all-select")
    }
}


/**
 * 弹出选择框
 * 2018-01-13 15:03
 * @param title
 * @param type
 * @param yes
 * @param no
 * @param okmsg
 * @param failmsg
 * @param func
 * @param reload
 */
function Swal(title,type,yes,no,okmsg, failmsg, func, reload) {
    if(!type){
        type  = "warning"
    }
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-template').click(function () {
            swal({
                title: title,
                text: "",
                type: type,
                showCancelButton: true,
                confirmButtonText: yes,
                cancelButtonText: no,
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                success("已发送请求");
                var result = eval(func);
                if (JSON.stringify(result).indexOf(okmsg) != -1) {
                    swal(
                        okmsg,
                        result,
                        'success'
                    )
                    eval(reload);
                    setTimeout(function () {
                        // 重新加载数据
                        eval(reload);
                    }, 5000)
                    setTimeout(function () {
                        // 重新加载数据
                        eval(reload);
                    }, 15000)
                } else {
                    swal(
                        failmsg,
                        result,
                        'error'
                    )
                }

            }, function (dismiss) {
                // dismiss can be 'cancel', 'overlay',
                // 'close', and 'timer'
            })
            // });
        },
            $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
    }(window.jQuery),

//initializing
        function ($) {
            "use strict";
            $.SweetAlert.init()
        }(window.jQuery);
}


/**
 * 选择checkbox的值
 * @return {string}
 * 2018-01-13 15:11
 */
function getCheckInput(selector) {
    var all = [];
    selector = "." + selector
    $(selector).each(function () {
        if($(this).is(":checked")){
            all.push($(this).val())
        }
    });
    return all.join(",");
}

/**
 * 保存后提示信息
 * 2018-01-13 19:02
 * @param result
 */
function saveMsg(result) {
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        success(result);
    }else{
        faild(result);
    }
}

/**
 * 检查选择的数据是否是单条
 * @return {boolean}
 */
function checkSignValue() {
    var value = getCheckInput("all");
    var v = value.split(",");
    if (v.length > 1) {
        faild("只能选择一项,不能同时选择多项");
        return false;
    }
    return true
}


// 将选择好的用户设置好
// 2018-01-20 13:53
function setSelectUser(users, id) {
    var html = "";
    var users = users.split(",");
    for (var i=0;i<users.length;i++){
        if (users[i]){
            html += "<option class='mul' value="+ users[i] + ">" + users[i] + "</option>\n";
        }
    }
    if(html){
        $('#undo_contact_group_redo_to').html(html);
    }
    if(id){
        $("#"+id).html(id);
    }
}

// 搜索用户
// 2018-01-20 13:24
function searchUser(val, id, dataId) {
    var temp = "";
    if(!dataId) {
        var data = $('#select_user_id').val();
    }else{
        var data = $('#'+dataId).html();
    }
    if(data){
        data = data.split("\n");
    }else{
        return;
    }
    for (var i = 0; i < data.length; i++) {
        if (val ){
            if (data[i].indexOf(val) != -1) {
                temp += data[i];
            }
        }else{
            temp += data[i];
        }
    }
    $('#' + id).html(temp);
}

/**
 * 2018-01-20 13:21
 * 设置用户数据
 */
function setUserData(id) {
    var url = "/api/users/name";
    var result = eval(get({}, url));
    var html = "";
    var select = "";
    for (var i = 0; i < result.length; i++) {
        html += "<option class='mul' value=" + result[i]["UserName"] + ">" + result[i]["UserName"] + "</option>\n";
    }
    if(!id){
        id = "select_user_id";
    }
    $('#'+id).val(select + html);
    $('#'+id).html(select + html);
}

// 修改状态使用
function changeStatus(obj) {
    var value = $(obj).val();
    if (value == 0) {
        $(obj).val(1);
    } else {
        $(obj).val(0);
    }
}


// 将选择好的组设置好
// 2018-01-21 9:53
function setSelectGroups(groups,id) {
    var url = "/api/groups/map";
    var gmap = get({}, url);
    var gmap = gmap["Data"];
    var html = ""
    var groups = groups.split(",");
    for (var i=0;i<groups.length;i++){
        if (groups[i]){
            if(gmap[groups[i]]){
                html += "<option class='mulgroup' value="+ groups[i] + ">" + gmap[groups[i]] + "</option>\n";
            }
        }
    }
    if(html){
        $('#undo_perm_group_redo_to').html(html);
    }
    if(id){
        $("#"+id).html(html);
    }
}

// 搜索组
// 2018-01-21 9:24
function searchGroups(val, id) {
    var temp = "";
    var data = $('#select_groups_id').val();
    if(data) {
        data = data.split("\n");
    }else{
        return;
    }
    for (var i = 0; i < data.length; i++) {
        if (val ){
            if (data[i].indexOf(val) != -1) {
                temp += data[i];
            }
        }else{
            temp += data[i];
        }
    }
    $('#' + id).html(temp);
}

/**
 * 2018-01-21 9:21
 * 设置组数据
 */
function setGroupsData(id, key) {
    if(!key){
        key = "GroupsId";
    }
    var url = "/api/groups/name";
    var result = eval(get({}, url));
    var html = "";
    var select = "";
    for (var i = 0; i < result.length; i++) {
        html += "<option class='mulgroup' value=" + result[i][key] + ">" + result[i]["GroupsName"] + "</option>\n";
    }
    if(!id){
        id = "select_groups_id";
    }
    $('#'+id).val(select + html);
    $('#'+id).html(select + html);
}


/**
 * 获取默认数据
 * @param v
 * @return {*}
 */
function getValue(v) {
    if(v){
        return v;
    }
    return "";
}

