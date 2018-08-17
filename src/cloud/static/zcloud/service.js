// 添加应用
function addService() {
    var url = "/application/service/add?ClusterName=" + getClusterName()
    window.location.href = url
}


/**
 * 删除应用
 * @param id
 * @return {*}
 */
function deleteService(id, force) {
    if (!checkSignValue()) {
        return;
    }
    if (!id) {
        var id = getCheckInput("all");
    }
    var url = "/api/service/" + id + "?force=" + force;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}


/**
 * 设置删除模板的id
 * @param id
 */
function setServiceDeleteId(id, force) {
    $("#delete_service_id").val(id)
    $("#delete_service_force_id").val(id)
    deleteServiceSwal();
}

/**
 * 到应用详情页面
 * @param name
 */
function toServiceDetail(name) {
    var url = "/application/service/detail/" + name + "?clusterName=" + getClusterName();
    window.location.href = url;
}

/**
 * 到容器详情页面
 * 2018-01-16 08:59
 * @param name
 */
function toContainerDetail(name) {
    var url = "/application/container/detail/" + name + "?clusterName=" + getClusterName();
    window.location.href = url;
}

/**
 * 获取服务标签的数据
 * @return {Array}
 */
function getLabelsData() {
    var labels = [];
    var count = 0;
    var d = {}
    $('.labelsdiv').each(function () {
        if (count == 0) {
            d["Key"] = $(this).html();
        }
        if (count == 1) {
            d["Value"] = $(this).html();
        }
        if (count == 2) {
            d["K8s"] = $(this).html();
            labels.push(d)
            count = 0;
        } else {
            count += 1;
        }
    });
    return JSON.stringify(labels);
}

/**
 * 获取健康检查的数据
 * @return {Array}
 */
function getConfigData() {
    var count = 0;
    var data = [];
    var d = {};
    $('.configuresdiv').each(function () {
        if (count == 0) {
            d["ContainerPath"] = $(this).html();
        }
        if (count == 1) {
            d["DataName"] = $(this).html();
        }
        if (count == 2) {
            d["DataId"] = $(this).html();
            data.push(d)
            d = {}
            count = 0;
        } else {
            count += 1;
        }
    });
    return JSON.stringify(data);
}

/**
 * 获取健康检查的数据
 * @return {Array}
 */
function getStorageData() {
    var storage = [];
    var count = 0;
    var d = {}
    $('.storagediv').each(function () {
        console.log($(this).html())
        if (count == 0) {
            d["ContainerPath"] = $(this).html();
        }
        if (count == 1) {
            var v = $(this).html();
            if (v != "undefined") {
                d["Volume"] = v;
            } else {
                d["Volume"] = "";
            }
        }
        if (count == 2) {
            d["HostPath"] = $(this).html();
            storage.push(d)
            d = {}
            count = 0;
        } else {
            count += 1;
        }
    });
    return JSON.stringify(storage);
}

/**
 * 获取健康检查的数据
 * @return {Array}
 */
function getHealthData() {
    var count = 0;
    var d = {};
    $('.healthdiv').each(function () {
        if (count == 0) {
            d["HealthType"] = $(this).html();
        }
        if (count == 1) {
            d["HealthPort"] = $(this).html();
        }
        if (count == 2) {
            if (d["HealthType"] == "CMD") {
                d["HealthCmd"] = $(this).html();
            } else {
                d["HealthPath"] = $(this).html();
            }
        }
        if (count == 3) {
            d["HealthInitialDelay"] = $(this).html();
        }
        if (count == 4) {
            d["HealthInterval"] = $(this).html();
        }
        if (count == 5) {
            d["HealthFailureThreshold"] = $(this).html();
        }
        if (count == 6) {
            d["HealthTimeout"] = $(this).html();
            count = 0;
        } else {
            count += 1;
        }
    });
    return d;
}
/**
 * 保存资源配额
 */
function saveService(serviceId) {
    if (!serviceId) {
        serviceId = 0;
    }

    $("#add_health_html").html("");
    $("#add_storage_html").html("");
    $("#add_configure_html").html("");
    var data = get_form_data();

    // 健康检查数据
    data["HealthData"] = JSON.stringify(getHealthData());
    data["StorageData"] = getStorageData();
    data["ServiceId"] = parseInt(serviceId);
    data["ServiceLablesData"] = getLabelsData();
    data["ConfigureData"] = getConfigData();
    data["ImageTag"] = data["ImageRegistry"] + ":" + data["Version"];
    console.log(data)
    if (!checkValue(data, "ServiceName,ImageRegistry,Version,ContainerPort,Replicas,ServiceName")) {
        return;
    }
    console.log(data)
    var url = "/api/service";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        $("#add_post_html").modal("toggle");
        success(result);
        setTimeout(function () {
            window.location.href = "/application/app/list";
        }, 3000);
        setTimeout(function () {
            window.location.href = "/application/app/list";
        }, 7000);
    } else {
        faild(result);
    }
}

// /**
//  * 删除模板弹出框
//  */
// function deleteServiceSwal() {
//     !function ($) {
//         "use strict";
//
//         var SweetAlert = function () {
//         };
//         //examples
//         SweetAlert.prototype.init = function () {
//             // //Parameter
//             // $('.delete-template').click(function () {
//             swal({
//                 title: '删除该应用',
//                 text: "",
//                 type: 'warning',
//                 showCancelButton: true,
//                 confirmButtonText: '确认删除',
//                 cancelButtonText: '不删除',
//                 confirmButtonClass: 'btn btn-success',
//                 cancelButtonClass: 'btn btn-danger m-l-10',
//                 buttonsStyling: false
//             }).then(function () {
//                 var forceId = $("#delete_service_force_id").val();
//                 var result = deleteService($("#delete_service_id").val(), forceId);
//                 if (result.indexOf("删除成功") != -1) {
//                     swal(
//                         '删除成功!',
//                         result,
//                         'success'
//                     );
//                     setTimeout(function () {
//                         loadServiceData()
//                     }, 2000);
//                 } else {
//                     swal(
//                         '删除失败!',
//                         result,
//                         'error'
//                     );
//                 }
//
//             }, function (dismiss) {
//                 // dismiss can be 'cancel', 'overlay',
//                 // 'close', and 'timer'
//                 $("#delete_service_id").val("")
//                 $("#delete_service_force_id").val("")
//             })
//             // });
//         },
//             $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
//     }(window.jQuery),
//
// //initializing
//         function ($) {
//             "use strict";
//             $.SweetAlert.init()
//         }(window.jQuery);
// }

/**
 *
 * @param obj
 */
function setBorderOut(obj) {
    obj.css("border", "1px solid #e0e0e0")
}

/**
 *
 * @param cpu
 * @param mem
 * @param id
 * @param custom
 */
function setSelectConfig(cpu, mem, id, custom) {
    if (custom) {
        $("#show-custom-config").show();
    } else {
        $("#show-custom-config").hide();
    }
    $(".service-header").css("background-color", "#ffffff");
    $(".showfa").hide();
    $("#" + id).css("background-color", "#eee");
    $("#" + id + "fa").show();
    $("input[name='Cpu']").val(cpu);
    $("input[name='Memory']").val(mem * 1024);
}


/**
 *
 * @param t
 */
function setServiceStatus(t) {
    if (t == 1) {
        $("#show_storage_id").show();
        $("#is_status").show();
        $("#no_status").hide();
        $("#is_status_border").css("cssText", "background-color:#eeeeee  !important;");
        $("#no_status_border").css("cssText", "background-color:#ffffff  !important;");
    } else {
        $("#show_storage_id").hide();
        $("#is_status").hide();
        $("#no_status").show();
        $("#no_status_border").css("cssText", "background-color:#eeeeee  !important;");
        $("#is_status_border").css("cssText", "background-color:#ffffff  !important;");
    }
}

/**
 *
 * @param id
 * @param value
 */
function setDeploy(id, value) {
    var http = "http";
    var tcp = "tcp";
    var cmd = "cmd";
    if (id == http || id == tcp || id == cmd) {
        $(".status").hide();
    } else {
        $(".statusd").hide();
    }
    $("#" + id).show();
    $(".health-check-select").css("cssText", "background-color:#ffffff  !important;");
    $(".service-mode-select").css("cssText", "background-color:#ffffff  !important;");
    $("#" + id + "_border").css("cssText", "background-color:#eeeeee  !important;");
    $("#statefulset_border").css("border-right", "1px solid #ccc");
    $("#cmd_border").css("border-right", "1px solid #ccc");
    if (id == http || id == tcp || id == cmd) {
        $("input[name='HealthType']").val(id.toUpperCase());
        if (id == tcp) {
            $("#tcp_hidden").hide();
        }
        if (id == http) {
            $("#tcp_hidden").show();
        }
        $(".cmd_hide").hide();
        $(".porthide").show();
    } else {
        $("input[name='serviceType']").val(value);
    }
    $("#is_status_border").show();
    if (id == "statefulset") {
        $("#is_status_border").hide();
    }
    if (id == cmd) {
        $("#tcp_hidden").hide();
        $(".porthide").hide();
        $(".cmd_hide").show();
    }
}


/**
 *
 * @param t
 */
function setUpdateMode(t) {
    if (t == 1) {
        $("#model_auto").show();
        $("#model_custom").hide();
        $("#update_model_auto").css("cssText", "background-color:#eeeeee  !important;");
        $("#update_model_custom").css("cssText", "background-color:#ffffff  !important;");
    } else {
        $("#model_auto").hide();
        $("#model_custom").show();
        $("#update_model_custom").css("cssText", "background-color:#eeeeee  !important;");
        $("#update_model_auto").css("cssText", "background-color:#ffffff  !important;");
    }
}

/**
 * 添加持久化存储弹出框
 */
function addStorage() {
    var url = "/application/service/storage/add";
    var result = get({ClusterName:$("#select-cluster-id").val()}, url);
    $("#add_health_html").html("")
    $("#add_storage_html").html(result);
    $("#add_post_html").modal("toggle");
}

/**
 * 添加健康检查弹出框
 */
function addHealth() {
    var url = "/application/service/health/add";
    var result = get({}, url);
    $("#add_storage_html").html("");
    $("#add_health_html").html(result);
    $("#add_post_html").modal("toggle");
}

/**
 * 添加修改端口弹出框
 * 2018-01-14 13:29
 */
function addPort() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/port/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}


/**
 *  添加自动扩容页面
 */
function addScale() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/scale/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}


/**
 * 服务修改健康检查弹出页面
 * 2018-01-14 12:29
 */
function changeHealth() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/health/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}

/**
 * 弹出服务修改cpu，内存实例
 * 2018-01-13 18:41
 */
function addCpuMemory() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/config/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}

/**
 * 弹出服务灰度升级的页面
 * 2018-01-14 09:28
 */
function addServiceImage() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/image/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}


/**
 * 弹出环境变量修改页面
 * 2018-01-14 11:15
 */
function addChangeEnv() {
    if (!checkSignValue()) {
        return
    }
    var id = getCheckInput("all")
    var url = "/application/service/env/add/" + parseInt(id);
    var result = get({}, url);
    $("#add_scale_html").html(result);
    $("#add_post_html").modal("toggle");
}


/**
 * 添加健康检查弹出框
 */
function addServiceConfig() {
    var cluster = $("#select-cluster-id").val();
    var entname = $("#select-entname-id").val();
    if(!cluster || ! entname || cluster.indexOf("请选择") != -1 || entname.indexOf("请选择") != -1){
        faild("请选择集群和环境");
        return
    }
    var url = "/application/service/configure/add";
    var result = get({}, url);
    $("#add_health_html").html("");
    $("#add_storage_html").html("");
    $("#add_configure_html").html("");
    $("#add_configure_html").html(result);
    $("#add_post_html").modal("toggle");
}

/**
 * 设置添加持久化数据
 */
function setConfigureData() {
    var c = $("input[name='ContainerPath']").val();
    var d = $("input[name='DataName']").val();
    var id = $("input[name='DataId']").val();
    var data = $("input[name='Data']").val();
    if(!c || !d||!data){
        return
    }

    if (id) {
        $("#" + id.replace(/DDATA/, "")).remove();
    }
    if (!checkUniqueData("configuresKey", c)) {
        setError($("input[name='ContainerPath']"), "uniquemsg");
        return;
    }
    var template = $("#configure_template_id").val();
    template = template.replace(/ContainerPath/g, c);
    template = template.replace(/DataId/g, data);
    template = template.replace(/DataName/g, d);
    template = template.replace(/TEMPLATEIDC/g, guid());
    $("#add_config_html_value").append(template);
}


/**
 * 设置添加持久化数据
 */
function setStorageData() {
    var c = $("input[name='ContainerPath']").val();
    var h = $("input[name='HostPath']").val();
    var v = $("select[name='Volume']").val();
    var id = $("input[name='StorageId']").val();
    console.log(v)
    if (id) {
        $("#" + id).remove();
    }
    if (!checkUniqueData("containerPath", c)) {
        return;
    }
    var template = $("#storage_template_id").val();
    template = template.replace(/CONTAINER_PATH/g, c);
    template = template.replace(/STORAGE_PATH/g, v);
    template = template.replace(/HOST_PATH/g, h);
    template = template.replace(/TEMPLATEID/g, guid());
    $("#add_storage_html_value").append(template);
}

/**
 * 获取uuid
 * @return {string}
 */
function guid() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

/**
 * 删除网页里某个元素的网页内容
 * @param id
 */
function removeData(id) {
    $("#" + id).remove();
}


/**
 * 检查数据是否有重复的
 * @param volumn
 * @return {boolean}
 */
function checkUniqueData(id, value) {
    var result = true;
    id = "." + id;
    $(id).each(function () {
        if (value == $(this).html()) {
            result = false
        }
        // ports.push($(this).html());
    });
    return result;
}

/**
 * 检查健康检查数据正确性
 * @param type
 * @param port
 * @param cmd
 * @return {boolean}
 * 2018-01-14 17:02
 */
function checkHealthData(type, port, cmd, path) {
    var data = {}
    if (type.toUpperCase() == "TCP" || type.toUpperCase() == "HTTP") {
        data["port"] = port
        if (!checkPort(data, "port")) {
            setInputError($("#HealthPort"), "errmsg")
            return false
        }
    }
    if (type.toUpperCase() == "HTTP") {
        if (!path || path.split("")[0] != "/") {
            setInputError($("#HealthPath"), "errmsg")
            return false
        }
    }
    if (type.toUpperCase() == "CMD") {
        data["cmd"] = cmd
        if (!checkValue(data, "cmd")) {
            setInputError($("#HealthCmd"), "errmsg")
            return false
        }
    }
    return true
}

/**
 * 渲染健康检查模板
 */
function setHealthData() {
    var type = $("input[name='HealthType']").val();
    var port = $("input[name='HealthPort']").val();
    var path = $("input[name='HealthPath']").val();
    var cmd = $("input[name='HealthCmd']").val();
    var delay = $("input[name='HealthInitialDelay']").val();
    var intervval = $("input[name='HealthInterval']").val();
    var failure = $("input[name='HealthFailureThreshold']").val();
    var timeout = $("input[name='HealthTimeout']").val();
    var data = {}
    if (type.toUpperCase() == "TCP" || type.toUpperCase() == "HTTP") {
        data["port"] = port
        if (!checkPort(data, "port")) {
            setInputError($("#HealthPort"), "errmsg")
            return
        }
    }
    if (type.toUpperCase() == "CMD") {
        data["cmd"] = cmd
        if (!checkValue(data, "cmd")) {
            setInputError($("#HealthCmd"), "errmsg")
            return
        }
    }

    var id = $("input[name='HealthId']").val();
    if (id) {
        $("#" + id).remove();
    }
    // 检查是否有重复的端口
    if ($("#add_health_html_value").html().indexOf("div") != -1) {
        return;
    }
    var template = $("#health_template_id").val();
    template = template.replace(/协议/g, type);
    if (type.toUpperCase() == "CMD") {
        template = template.replace(/端口/g, "");
        template = template.replace(/路径/g, cmd);
    } else {
        template = template.replace(/路径/g, path);
        template = template.replace(/端口/g, port);
    }
    if (type.toUpperCase() == "TCP") {
        template = template.replace(/路径/g, "");
    }
    template = template.replace(/启动预估/g, delay);
    template = template.replace(/间隔/g, intervval);
    template = template.replace(/超时/g, timeout);
    template = template.replace(/不健康阈值/g, failure);
    template = template.replace(/TEMPLATEIDH/g, guid());
    $("#add_health_html_value").append(template);
}


/**
 * 获取配置文件的数据
 * @param id
 * @return {Array|*}
 */
function getConfigureValue(id) {
    var value = $("#" + id).html()
    value = value.replace(/ /g, "");
    value = value.replace(/<\/div>/g, ",");
    value = value.replace(/<divclass="col-md-5fw100control-labelconfiguresdivconfiguresKey">/g, "");
    value = value.replace(/<divclass="col-md-4fw100control-labelconfiguresdiv">/g, "");
    value = value.replace(/<divclass="form-group"><divclass="col-md-12"style="margin-left:-25px;"><divclass="form-group">/g, "");
    value = value.replace(/\n/g, "");
    var values = value.split(",");
    return values
}

/**
 * 编辑健康检查
 * 20180109
 * @param id
 */
function editHealthData(id) {
    addHealth();
    var values = getHealthData();
    var type = values["HealthType"].toLowerCase();
    console.log(values)
    setDeploy(values["HealthType"].toLowerCase());
    $("input[name='HealthType']").val(values["HealthType"]);
    $("input[name='HealthPort']").val(values["HealthPort"]);
    if (type == "cmd") {
        $(".cmd_hide").show();
        $("input[name='HealthCmd']").val(values["HealthPath"]);
    }
    if (type == "http") {
        $("#tcp_hidden").show();
        $(".healthpath").show();
        $("input[name='HealthPath']").val(values["HealthPath"]);
    }
    if (type == "tcp") {

        $("#tcp_hidden").hide();
    }
    $("input[name='HealthInitialDelay']").val(values["HealthInitialDelay"]);
    $("input[name='HealthInterval']").val(values["HealthInterval"]);
    $("input[name='HealthFailureThreshold']").val(values["HealthFailureThreshold"]);
    $("input[name='HealthTimeout']").val(values["HealthTimeout"]);
    $("input[name='HealthId']").val(id)
}

/**
 * 编辑存储数据
 * @param id
 */
function editStorageData(id) {
    addStorage();
    var value = $("#" + id).html();
    value = value.replace(/ /g, "");
    value = value.replace(/<\/div>/g, ",");
    value = value.replace(/<divclass="col-md-4fw100control-labelwrapstoragediv">/, "");
    value = value.replace(/<divclass="col-md-4fw100control-labelwrapstoragedivcontainerPath">/, "");
    value = value.replace(/<divclass="col-md-2fw100control-label">/g, "");
    value = value.replace(/<divclass="col-md-4fw100control-labelwrap">/g, "");
    value = value.replace(/\n/g, "");
    var values = value.split(",");
    $("input[name='ContainerPath']").val(values[0]);
    $("input[name='HostPath']").val(values[2]);
    $("input[name='Volume']").val(values[1]);
    $("input[name='StorageId']").val(id);
}

/**
 * 添加服务标签提示错误信息
 * @param obj1
 * @param obj2
 */
function setError(obj1, obj2) {
    if (obj1) {
        setInputError(obj1, "errmsg");
    }
    if (obj2) {
        setInputError(obj2, "errmsg");
    }
}

/**
 * 添加服务标签
 */
function addLabels() {
    //
    var key = $("input[name='LabelsKey']").val();
    var value = $("input[name='LabelsValue']").val();
    var k8s = $("input[name='LabelsK8s']").is(":checked");
    var objkey = $("input[name='LabelsKey']");
    var objvalue = $("input[name='LabelsValue']");
    if (!key || !value) {
        setError(objkey, objvalue);
        return;
    } else {
        if (key.length < 2 || value.length < 1) {
            setError(objkey, objvalue);
            return
        }
        if (!checkClusterName(key)) {
            setError(objkey, objkey);
            return
        }
        setInputOk(objvalue);
        setInputOk(objkey);
    }
    if (!checkUniqueData("lableskey", key)) {
        setInputError(objkey, "uniquemsg");
        return;
    }
    var template = $("#labels_template_id").val();
    template = template.replace(/键/g, key);
    template = template.replace(/值/g, value);
    template = template.replace(/Kubernetes/g, k8s);
    $("#add_labels_html_value").append(template);
}


/**
 * 加载集群在页面的集群选择信息
 * 2018-01-04
 */
function loadConfigureSelect() {
    var url = "/api/configure/name";
    var html = "<option>--请选择--</option>";
    var cluster = $("#select-cluster-id").val();
    var entname = $("#select-entname-id").val();
    if(!cluster || ! entname || cluster.indexOf("请选择") != -1 || entname.indexOf("请选择") != -1){
        faild("请选择集群和环境");
        return
    }
    var result = get({ClusterName:cluster, Entname:entname}, url);
    //  先获取cookie选择好的
    for (var i = 0; i < result.length; i++) {
        html += "<option value='" + result[i]["ConfigureId"] + "'>" + result[i]["ConfigureName"] + "</option>"
    }
    var id = "add_config_select_id";
    $("#" +id ).html(html);
    $("#" + id).selectpicker('refresh');
}

/**
 * 编辑配置文件数据
 */
function editConfigureData(id) {
    addServiceConfig();
    var values = getConfigureValue(id);
    $("input[name='ContainerPath']").val(values[0]);
    $("input[name='ConfigureName']").val(values[1]);
    $("input[name='DataName']").val(values[2]);
    $("input[name='DataId']").val(id)
}

/**
 * 显示配置文件内容
 */
function showConfigure(id) {
    var url = "/api/configure/data/" + id;
    var result = get({}, url);
    result = result["data"];
    $("#show_configure_id").val(result[0]["Data"]);
    $("#show_configure_data_group").val(result[0]["ConfigureName"])
    $("#show_configure_data_name").val(result[0]["DataName"])
    $("#show_configure_id_pop").modal("toggle");
}

function loadConfigItem(configId) {
    if (!configId) {
        return;
    }
    $("input[name='Data']").val("");
    var url = "/api/configure/data?ConfigureId=" + configId;
    var html = "";
    var result = get({}, url);
    result = result["data"];
    //  先获取cookie选择好的
    for (var i = 0; i < result.length; i++) {
        html += "<option value='" + result[i]["DataName"] + "'>" + result[i]["DataName"] + "</option>"
    }
    $("#select_data_name_id").html(html);
    $('.selectpicker').selectpicker('refresh');
}


// 将服务全部关闭,量为0
function stopServiceSwal() {
    Swal("将停止所有选择的服务", "warning", "确认停止", "不停止", "成功", "失败", "stopService()", "loadServiceData()");
}


/**
 * 停止服务
 * 2018-01-13 15:18
 * @return {*}
 */
function stopService() {
    if (!checkSignValue()) {
        return
    }
    var value = getCheckInput("all");
    var url = "/api/service/scale/" + value;
    var result = get({replicas: 0}, url);
    return JSON.stringify(result);
}


// 将服务全部关闭,并重新启动
function restartServiceSwal() {
    Swal("将快速重启该服务", "warning", "确认重启", "不重启", "重启成功", "重启失败", "restartService()", "loadServiceData()");
}


/**
 * 启动
 * 2018-01-13 15:37
 * @return {*}
 */
function startService() {
    var value = getCheckInput("all");
    var url = "/api/service/scale/" + parseInt(value);
    var result = get({replicas: 0, "start": 1}, url);
    saveMsg(result);
    setTimeout(function () {
        loadServiceData()
    }, 2000)
}

/**
 * 2018-01-13 17:43
 * 水平扩展保存
 */
function serviceScale() {
    var url = "/api/service/scale/" + parseInt($("#scale_service_id").val());
    var result = get({replicas: parseInt($("#range_01").val())}, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1) {
        success(result);
        $("#add_post_html").modal("toggle");
        $("#add_scale_html").html("");
    } else {
        faild(result);
    }
}

/**
 * 2018-01-13 15:58
 * 重启服务
 */
function restartService() {
    var value = getCheckInput("all");
    var url = "/api/service/scale/" + parseInt(value);
    get({replicas: 0}, url);
    var result = get({replicas: 0, "start": 1}, url);
    saveMsg(result);
}

/**
 * 服务配置修改
 * 2018-01-13 19:23
 */
function serviceCpuSave() {
    var cpu = $("input[name='Cpu']").val();
    var mem = $("input[name='Memory']").val();
    var id = $("#scale_service_id").val();
    var url = "/api/service/update/" + id
    var result = post({cpu: cpu, mem: mem, type: "config"}, url)
    saveMsg(result);
}

/**
 * 升级镜像
 * 2018-01-14 08:39
 */
function updateServiceImage() {
    var id = $("#update_image_service_id").val();
    var version = $("#update_image_version_id").val();
    var interval = $("#update_interval_id").val();
    if (!version) {
        setInputError($("#update_image_version_id"), "errmsg")
        return
    }
    if (!interval) {
        setInputError($("#update_interval_id"), "errmsg")
        return
    }
    var url = "/api/service/update/" + id
    var result = post({version: version, type: "image", MinReady: interval}, url)
    saveMsg(result);
}

/**
 * 服务环境变量修改
 * 2018-01-14 11:26
 */
function envSave() {
    var env = $("#env_change_id").val();
    var id = $("#update_service_env_id").val();
    var url = "/api/service/update/" + id;
    var result = post({env: env, type: "env"}, url);
    saveMsg(result);
}

/**
 * 保存健康检查数据
 * 2018-01-14 12:02
 */
function saveHealth() {
    var d = {}
    d["HealthType"] = $("input[name='HealthType']").val();
    d["HealthPort"] = $("input[name='HealthPort']").val();
    d["HealthCmd"] = $("input[name='HealthCmd']").val();
    d["HealthPath"] = $("input[name='HealthPath']").val();
    d["HealthInterval"] = $("input[name='HealthInterval']").val();
    d["HealthFailureThreshold"] = $("input[name='HealthFailureThreshold']").val();
    d["HealthTimeout"] = $("input[name='HealthTimeout']").val();
    if (!checkHealthData(d["HealthType"], d["HealthPort"], d["HealthCmd"], d["HealthPath"])) {
        return
    }
    var id = $("#update_service_health_id").val();
    var url = "/api/service/update/" + id;
    var result = post({healthData: JSON.stringify(d), type: "health"}, url);
    saveMsg(result);
}

/**
 * 保存端口数据
 * 2018-01-14 13:27
 */
function portSave() {
    var ports = $("#ports_change_id").val();
    var id = $("#update_service_port_id").val();
    var url = "/api/service/update/" + id;
    var result = post({port: ports, type: "port"}, url);
    saveMsg(result);
}

// 2018-01-15 15:06
// 连接到docker容器
function toTty(id) {
    window.open("/webtty/" + id)
}


// 容器删除后跳到容器列表
// 2018-01-16 12:47
function toContainerList() {
    setTimeout(function () {
        window.location.href = "/application/container/list";
    }, 3000);
}

// 删除容器
function deleteServiceSwal(id, force) {
    id = getValue(id);
    force = getValue(force);
    var msg = "将删除服务";
    if (force) {
        msg = "强制删除该服务,并不会检查k8s集群可用性";
    }
    Swal(msg, "warning", "确认操作", "不操作", "成功", "失败", "deleteService('" + id + "','" + force + "')", "loadServiceData()");
}

// 删除容器
function deleteContainerSwal(id) {
    Swal("将停止或删除该容器", "warning", "确认操作", "不操作", "成功", "失败", "deleteContainer(" + id + ")", "toContainerList()");
}

// 删除容器
// 2018-01-16 12:43
function deleteContainer(id) {

    if (!id) {
        var result;
        var ids = getCheckInput("all").split(",")
        if (ids.length < 1) {
            return "至少选择一项";
        }
        for (var i = 0; i < ids.length; i++) {
            var url = "/api/container/" + parseInt(ids[i]);
            result = del({}, url);
        }
        return JSON.stringify(result)
    } else {
        var url = "/api/container/" + parseInt(id);
        var result = del({}, url);
        return result;
    }
}


function loadServiceData(key,name) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }
    name = getValue(name);
    // 应用详情页面使用 2018-01-18 08:02
    var app = $("#detail_add_app_name").val();
    if (!app) {
        app = ""
    }

    $("#service-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bLengthChange": false,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/service?t=" + new Date().getTime() + "&key=" + key + "&AppName=" + app +"&ServiceName="+name,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "ServiceId", "sWidth": "4%", "mRender": function (data) {
                var html = '<div class="checkbox checkbox-primary ">' +
                    '<input id="checkbox0" onchange="checkBoxChange($(this))" type="checkbox" class="all" value="' + data + '">' +
                    '<label for="checkbox0">' +
                    '</label>' +
                    '</div>';
                html = html.replace(/checkbox0/g, "checkbox" + data);
                return html;
            }
            },
            {
                "data": "ServiceName", "sWidth": "8%", "mRender": function (data, type, full) {
                return "<a href='javascript:void(0)' onclick='toServiceDetail(\"" + data + "\")'>" + data + "<br></a><span class='text-default'>应用:&nbsp;" + full["AppName"] + "</span>";
            }
            },
            {
                "data": "Status", "sWidth": "8%", "mRender": function (data, type, full) {
                if (data == "True") {
                    var r = '<div class="Running"><div><i class="fa fa-circle"></i><span>&nbsp;运行中</span>' +
                        "<div class='text-default'>" + full["AvailableReplicas"] + "/" + full["ContainerNumber"] + "&nbsp;共" + full["AvailableReplicas"] + "个运行</div>"
                        + '</div></div>'
                }
                if (full["AvailableReplicas"] == 0) {
                    var r = '<div class="Fail"><div><i class="fa fa-circle"></i><span>&nbsp;已停止</span>&nbsp;' +
                        "<div class='text-default'>" + full["AvailableReplicas"] + "/" + full["ContainerNumber"] + "&nbsp;共" + full["AvailableReplicas"] + "个运行</div>"
                        + '</div>' +
                        '</div></div>'
                }
                if (!r) {
                    return "<span class='Fail'>未知</span>"
                }
                var errmsg = full["ErrorMsg"];
                if (errmsg) {
                    return "<span title='" + full["ErrorMsg"] + "'>" + r + "</span>"
                } else {
                    return "<span>" + r + "</span>"
                }
            }
            },
            {
                "data": "Image", "sWidth": "12%", "mRender": function (data) {
                return "<div style='word-wrap:break-word'><a>" + data + "</a></div>";
            }
            },
            {
                "data": "ResourceName", "sWidth": "13%", "mRender": function (data, type, full) {
                return "<a href='/base/quota/detail/" + data + "'>" + data + "</a><br>集群名称:&nbsp;" + full["ClusterName"]
            }
            },
            {"data": "CreateTime", "sWidth": "7%"},
            {
                "data": "Access", "sWidth": "16%", "mRender": function (data,type, full) {

                if (data) {
                    data = data.join("<br>")
                }
                    if(full["Domain"]){
                        data += "<br>域名访问: <a target='_blank' href='http://" + full["Domain"] +"'/>"+full["Domain"]+"</a><br>"
                    }
                   if(data){
                    return data;
                   }
                return "<span class='Fail'>未知</span>"

            }
            },
            {
                "data": "ServiceId", "sWidth": "5%", "mRender": function (data) {
                return '<button type="button"  title="强制删除" onClick="deleteServiceSwal(' + data + ',1)" class="delete-groups btn btn-xs rb-btn-oper"><i class="fa fa-bolt"></i></button>';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
    $("#app-data-table_wrapper").css("cssText", "margin-top:-20px !important;");

}

/**
 * 2018-02-07 13:58
 * 镜像选择后设置镜像名称
 */
function setImageName(name,obj) {
    $("input[name='ImageRegistry']").val(name);
    var temp = name.split("/");
    var name = [];
    for(var i=1;i<temp.length;i++){
        name.push(temp[i]);
    }
    $("#select-registry-group").val(temp[1]);
    getImageTags(name.join("/"));
    $(".btnall").css("background-color","#4cb7e8");
    $(obj).css("background-color", "#ffa91c");
    // $('.button-next.btn.btn-primary.waves-effect.waves-light').trigger("click");
}

/**
 * 2018-02-07 11:14
 * 加载镜像部署数据
 * @param key
 * @param grouptype
 */
function loadImageData(key, clusterName) {
    if (!key) {
        key = $("#search_user_id").val();
    } else {
        if (key.length < 4) {
            return
        }
    }
    key = getValue(key);
    clusterName = getValue(clusterName);
    $("#load-image-table").dataTable({
        "filter": true,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "bLengthChange": false,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": false, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/registry/deploy/images?t=" + new Date().getTime() + "&search=" + key + "&clusterName=" + clusterName,
            "type": 'get'
        },
        "columns": [ // 数据映射
            {
                "data": "ClusterName", "sWidth": "20%", "mRender": function (data) {
                return data;

            }
            },
            {
                "data": "ServerDomain", "sWidth": "60%", "mRender": function (data, type, full) {
                return data + "/" + full["Name"];
            }
            },
            {
                "data": "ServerDomain", "sWidth": "20%", "mRender": function (data,type,full) {
                return '<button type="button" style="padding: 5px !important;"  title="更新" onclick="setImageName(\'' + data + "/" + full["Name"] + '\',$(this));" class="btnall  btn  rb-btn-oper pull-right m-r-10"><i class="fa fa-arrow-circle-o-right"></i>&nbsp;使用该镜像部署</button>&nbsp;';
            }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}
/*
 /**
 * Created by zhaoyun on 2018/1/5.
 */

/**
 * 获取镜像
 * 2018-02-06 14:51
 * @param id
 */
function getImageTags(id) {
    var url = "/api/registry/group/images/" + id;
    var result = get({GroupName: $("#select-registry-group").val()}, url);
    var data = result["data"];
    var html = "";
    var tags = data["Tags"];
    tags = tags.split(",");
    var newTag = [];
    for (var i = tags.length; i >= 0; i--) {
        if (tags[i]) {
            newTag.push(tags[i]);
        }
    }
    for (var i = 0; i < newTag.length; i++) {
        html += "<option value='" + newTag[i] + "'>" + newTag[i] + "</option>"
    }
    console.log(html)
    $("#select-version-id").html(html);
    $("#select-version-id").selectpicker('refresh');
}
