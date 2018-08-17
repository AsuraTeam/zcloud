
// 添加模板
function addTemplate(templateId) {
    if(!templateId){
        templateId = 0
    }
    var url = "/application/template/add";
    var result = post({ClusterName: "{{.data.ClusterName}}", TemplateId:templateId}, url);
    $("#add_template_html").html(result)
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-08-16 09:57
 * 添加模板更新
 * @param templateId
 */
function addTemplateUpdate(templateId) {
    var url = "/application/template/update/add";
    var result = post({TemplateId:templateId}, url);
    $("#add_template_html").html(result)
    $("#add_post_html").modal("toggle")
}

/**
 * 2018-08-16 11:02
 * 应用拉起页面
 * @param templateId
 */
function addTemplateDeploy(templateId) {
    var url = "/application/template/deploy/add";
    var result = post({TemplateId:templateId}, url);
    $("#add_template_html").html(result)
    $("#add_post_html").modal("toggle")
}


/**
 * 设置删除模板的id
 * @param id
 */
function setDeleteId(id) {
    $("#delete_template_id").val(id);
    deleteTemplateSwal();
}

/**
 * 删除模板弹出框
 */
function deleteTemplateSwal() {
    !function ($) {
        "use strict";

        var SweetAlert = function () {
        };
        //examples
        SweetAlert.prototype.init = function () {
            // //Parameter
            // $('.delete-template').click(function () {
                swal({
                    title: '删除该模板',
                    text: "",
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonText: '确认删除',
                    cancelButtonText: '不删除',
                    confirmButtonClass: 'btn btn-success',
                    cancelButtonClass: 'btn btn-danger m-l-10',
                    buttonsStyling: false
                }).then(function () {
                    var result = deleteTemplate($("#delete_template_id").val())
                    if (result.indexOf("删除成功") != -1) {
                        swal(
                            '删除成功!',
                            result,
                            'success'
                        )
                        setTimeout(function () {
                            window.location.href = "/application/template/list"
                        },2000)
                    } else {
                        swal(
                            '删除失败!',
                            result,
                            'error'
                        )
                    }

                }, function (dismiss) {
                    // dismiss can be 'cancel', 'overlay',
                    // 'close', and 'timer'
                    $("#delete_template_id").val("")
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
 * 加载数据
 * @param key
 */
function loadHistoryData(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#template-data-h-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": true, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/template/deploy/history?t=" + new Date().getTime() + "&search=" + key + "&cluster="+getClusterName(),
            "type": 'post'
        },
        "columns": [ // 数据映射
            {"data": "TemplateName"},
            {"data": "Entname"},
            {"data": "ClusterName"},
            {"data": "ResourceName"},
            {"data": "AppName"},
            {"data": "ServiceName","mRender":function (data) {
                    return data;
                }},
            {"data": "Domain"},
            {"data": "CreateUser"},
            {"data": "CreateTime"},
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

loadTemplateData();


/**
 * 加载数据
 * @param key
 */
function loadTemplateData(key) {
    if (!key) {
        key = ""
    } else {
        if (key.length < 4) {
            return
        }
    }

    $("#template-data-table").dataTable({
        "filter": false,//去掉搜索框
        "ordering": false, // 是否允许排序
        "paginationType": "full_numbers", // 页码类型
        "destroy": true,
        "processing": true,
        "bPaginate": false, //是否显示（应用）分页器
        "serverSide": true,
        "bInfo": true, //是否显示页脚信息，DataTables插件左下角显示记录数
        "scrollX": true, // 是否允许左右滑动
        "displayLength": 10, // 默认长度
        "ajax": { // 请求地址
            "url": "/api/template?t=" + new Date().getTime() + "&search=" + key + "&cluster="+getClusterName(),
            "type": 'get'
        },
        "columns": [ // 数据映射
            {"data": "TemplateName"},
            {"data": "ServiceName","mRender":function (data) {
                    return data.split(",").length;
                }},
            {"data": "ServiceName","mRender":function (data) {
                return data.replace(/,/, "<br>");
            }},
            {"data": "Description"},
            {"data": "CreateTime"},
            {"data": "LastModifyTime"},
            {
                "sWidth": "150px", "data": "TemplateId", "mRender": function (data) {
                       return '<button type="button" title="编辑数据文件" onclick="addTemplateUpdate(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-edit"></i></button>&nbsp;' +
                        '<button type="button"  title="删除" onClick="setDeleteId(' + data + ')" class="delete-template btn btn-xs rb-btn-oper"><i class="fa fa-trash-o"></i></button>';
            }
            },
            {
                "sWidth": "150px", "data": "TemplateId", "mRender": function (data) {
                    return '<button type="button" title="应用拉起" onclick="addTemplateDeploy(' + data + ')" class="btn btn-xs rb-btn-oper"><i class="fa fa-send-o"></i></button>&nbsp;' ;
                }
            },
        ],
        "fnRowCallback": function (row, data) { // 每行创建完毕的回调
            $(row).data('recordId', data.recordId);
        }
    });
}

loadTemplateData();

/**
 * 删除模板方法
 * @param id
 * @return {*}
 */
function deleteTemplate(id) {
    var url = "/api/template/"+id
    var result = del({}, url)
    result = JSON.stringify(result)
    return result
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
 * 保存模板
 */
function saveTemplate(templateId) {
    if(!templateId){
        templateId = 0
    }
    var data = get_form_data();
    data["TemplateId"] = parseInt(templateId);
    if(!checkValue(data,"TemplateName")){
        return
    }
    var service = [];
    $("#undo_contact_group_redo_to option").each(function () {
        service.push($(this).val());
    });
    data["ServiceName"] = service.join(",");
    if(!checkValue(data,"ServiceName")){
        return
    }
    var url = "/api/template";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadTemplateData()
    }else{
        faild(result)
    }
}

/**
 * 2018-08-16 09:51
 * 保存模板更新
 */
function saveTemplateUpdate(templateId) {
    var data = get_form_data();
    data["TemplateId"] = parseInt(templateId);
    if(!checkValue(data,"TemplateName,Yaml")){
        return
    }
    var url = "/api/template/update";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadTemplateData()
    }else{
        faild(result)
    }
}

/**
 * 2018-08-16 10:51
 * 拉起模板
 */
function saveStartDeploy(templateId) {
    var data = get_form_data();
    data["TemplateId"] = parseInt(templateId);
    if(!checkValue(data,"TemplateName,AppName,Ent,Cluster,ResourceName")){
        return
    }
    var url = "/api/template/deploy/"+ templateId;
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle");
        success(result);
        loadTemplateData()
    }else{
        faild(result)
    }
}