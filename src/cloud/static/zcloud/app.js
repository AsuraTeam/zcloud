
// 添加应用
function addApp() {
   var appId = $("#detail_add_app_id").val();
   if(!appId){
       appId = 0;
   }
    var url = "/application/app/add?ClusterName="+getClusterName()+"&AppId="+appId;
    window.location.href = url;
    // var result = get({ClusterName: getClusterName(), AppId:appId}, url)
    // $("#add_app_html").html(result)
    // $("#add_post_html").modal("toggle")
}

function showTemplate() {
    if($("#check_template_show").is(":checked")){
        $("#show_templateid").show()
    }else{
        $("#show_templateid").hide()
    }
}

/**
 * 删除应用
 * 2018-01-17 08:47
 * @return {*}
 */
function deleteApp(id,force) {
    if(!checkSignValue() && !id){
        return
    }
    if (!id){
        id = getCheckInput("all")
    }
    if(!force){
        force = ""
    }
    var url = "/api/app/"+id
    var result = del({force:force}, url)
    result = JSON.stringify(result)
    return result
}


/**
 * 设置删除模板的id
 * @param id
 */
function setAppDeleteId(id,force) {
    $("#delete_app_id").val(id)
    $("#delete_app_force_id").val(id)
    deleteAppSwal();
}

/**
 * 到应用详情页面
 * @param name
 */
function toAppDetail(name,yaml) {
    var url = "/application/app/detail/" +name;
    if(yaml){
        url = url + "&yaml=1"
    }

    window.location.href = url;
}

/**
 * yaml  部署模板
 */
function setTemplate() {
    if (!$("textarea[name='Yaml']").val()) {
        $("textarea[name='Yaml']").val($("#yaml_template").val());
        checkChange('Yaml','textarea');
    }
}

/**
 * 添加应用使用模板方式
 * @param id
 */
function setYaml(id) {
    if(!id){
        return
    }
    var url = "/api/template/"+id
    var data = get({}, url)
    var yaml = data["data"][0]["Yaml"]
    $("#field-5").val(yaml)
}

/**
 * 保存资源配额
 */
function saveApp(appId) {
    if(!appId){
        appId = 0
    }
    var data = get_form_data();
    data["AppId"] = parseInt(appId)
    if(!checkValue(data,"AppName,Yaml,ResourceName")){
        return
    }
    if(!data["ResourceName"]){
        setInputError($("input[name='ResourceName']"), "nullmsg")
        return
    }
    var url = "/api/app";
    var result = post(data, url)
    result = JSON.stringify(result)
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle")
        success(result)
        loadAppData()
    }else{
        faild(result)
    }
}


/**
 * 启动
 * 2018-01-13 15:37
 * @return {*}
 */
function startApp() {
    if (!checkSignValue()){
        return
    }
    var value = getCheckInput("all");
    var url = "/api/app/scale/" + parseInt(value);
    var result = get({replicas: 0, "start": "1"}, url);
    saveMsg(result);
    setTimeout(function () {
        loadAppData();
    }, 2000)
}

/**
 * 停止拥有
 * 2018-01-16 21:17
 * @return {*}
 */
function stopApp() {
    if (!checkSignValue()){
        return
    }
    var value = getCheckInput("all");
    var url = "/api/app/scale/" + parseInt(value);
    var result = get({replicas: 0, "start": "0"}, url);
    return JSON.stringify(result)
}


// 停止应用
// 2018-01-16 21:19
function stopAppSwal() {
    Swal("将停止该应用", "warning", "确认操作", "不操作", "成功", "失败", "stopApp()", "loadAppData()");
}

// 重启应用
// 2018-01-16 21:21
function restartApp() {
    stopApp();
   return startApp();
}

/**
 * 重建应用
 * 新集群重建
 * 或者应用重新部署
 * 2018-02-26 09;17
 */
function redeploymentApp() {
        var value = getCheckInput("all");
        var url = "/api/app/redeploy";
        var result = get({apps:value}, url);
        // saveMsg(result);
        setTimeout(function () {
            loadAppData();
        }, 4000);
       return JSON.stringify(result);
}

// 重启应用
// 2018-01-16 21:20
function restartAppSwal() {
    Swal("将重启该应用", "warning", "确认操作", "不操作", "成功", "失败", "restartApp()", "loadAppData()");
}

// 重建应用
// 2018-08-13 08:50
function redeployAppSwal() {
    Swal("将重建该应用,重建是应用在集群失败或意外手动删除时,重新生成应用,许保证仓库是正常的", "warning", "确认操作", "不操作", "成功", "失败", "redeploymentApp()", "loadAppData()");
}

/**
 * 删除模板弹出框
 */
function deleteAppSwal(id,force) {
    if(!id){
        id = ""
    }
    if(!force){
        force = ""
    }
    Swal("将删除该应用", "warning", "确认操作", "不操作", "成功", "失败", " deleteApp(\'"+id+"\',\'"+force+"\')", "loadAppData()");
}


/**
 * 2018-02-06 14:34
 * 获取仓库组
 * @param registryName
 */
function selectImageGroup(registryName, target) {
    var url = "/api/registry/group";
    var result = get({ClusterName: $("#select-cluster-id").val(), ServerDomain: registryName}, url);
    var data = result["data"];
    var html = "<option>--请选择--</option>";
    for (var i = 0; i < data.length; i++) {
        html += "<option value='" + data[i]["GroupName"] + "'>" + data[i]["GroupName"] + "</option>"
    }
    var id = "select-registry-group";
    if (target) {
        id = "";
    }
    $("#" + id).html(html);
    $("#" + id).selectpicker('refresh');
}


/*
/**
 * Created by zhaoyun on 2018/1/5.
 */
