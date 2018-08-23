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

/**
 * 2018-08-23 09:41
 * 获取环境名称
 * @param change
 */
function setEntname() {
    var url = "/api/ent/name";
    var result = get({}, url);
    var html = "{{.entname}}<option>--请选择--</option>";
    for (var i = 0; i < result.length; i++) {
        var t = result[i]["Entname"];
        if (html.indexOf("'" + t) == -1 ) {
            html += "<option value='" + t + "'>" + t + "</option>";
        }
    }
    $("#select-entname-id").html(html);
    $('.selectpicker').selectpicker('refresh');
}

/**
 * 2018-01-27 18:18
 * @param cluster
 */
function setRegistryServer(cluster, target) {
    var url = "/api/registry";
    var result = get({ClusterName: cluster}, url);
    var data = result["data"];
    var html = "<option>--请选择--</option>";
    for (var i = 0; i < data.length; i++) {
        html += "<option value='"+data[i]["Name"]+"'>"+data[i]["Name"]+"</option>"
    }
    var id = "select-registry-server";
    if (target){
        id = "select-registry-server-target";
    }
    $("#"+id).html(html);
    $("#"+id).selectpicker('refresh');
}


/**
 * 2018-02-06 15:40
 * 获取集群名称
 * */
function getEntClusterData(entname, target) {
    var url = "/api/ent";
    var result = get({Entname: entname}, url);
    var data = result["data"];
    var html = "<option>--请选择--</option>";
    for (var i = 0; i < data.length; i++) {
        if (data[i]["Entname"] == entname) {
            var clusters = data[i]["Clusters"];
            var clusters = clusters.split(",");

            for (var j = 0; j < clusters.length; j++) {
                html += "<option value='" + clusters[j] + "'>" + clusters[j] + "</option>"
            }
        }
    }
    var id = "select-cluster-id";
    if (target) {
        id = "select-cluster-target-id";
    }
    var obj = $("#" + id);
    obj.html(html);
    obj.selectpicker('refresh');
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

/**
 * 获取镜像组下面的镜像
 * @param groupName
 */
function getImageName(groupName) {
    var url = "/api/registry/group/images";
    var result = get({
        ClusterName: $("#select-cluster-id").val(),
        ServerDomain: $("#select-registry-group").val(),
        GroupName: groupName
    }, url);
    var data = result["data"];
    var html = "<option>--请选择--</option>";
    for (var i = 0; i < data.length; i++) {
        html += "<option value='" + data[i]["Name"] + "'>" + data[i]["Name"] + "</option>"
    }
    $("#select-image-id").html(html);
    $("#select-image-id").selectpicker('refresh');
}

/**
 * 2018-02-03 21:02
 * 选择应用
 * @param clustername
 */
function getResouceData(clustername, id, url, name, appname) {
    appname = getValue(appname);
    var data = get({ClusterName: clustername, AppName: appname}, url);
    var html = "";
    for (var i = 0; i <= data.length; i++) {
        if (data[i]) {
            html += "<option value='" + data[i][name] + "'>" + data[i][name] + "</option>";
        }
    }
    $("#" + id).html(html);
    $('.selectpicker').selectpicker('refresh');
}


/**
 * 2018-02-03 21:02
 * 选择应用
 * @param clustername
 */
function getServiceIdSelect(clustername, id, url, name, appname) {
    appname = getValue(appname);
    var data = get({ClusterName: clustername, AppName: appname}, url);
    var html = "";
    for (var i = 0; i <= data.length; i++) {
        if (data[i]) {
            html += "<option value='" + data[i]["ServiceId"] + "'>" + data[i][name] + "</option>";
        }
    }
    $("#" + id).html(html);
    $('.selectpicker').selectpicker('refresh');
}
/**
 * Created by zhaoyun on 2018/2/7.
 */
