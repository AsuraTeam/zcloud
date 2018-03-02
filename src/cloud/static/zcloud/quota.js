
// 添加直接
function addQuota(quotaId) {
    if(!quotaId){
        quotaId = 0
    }
    var url = "/base/quota/add";
    var result = post({QuotaName: "{{.data.QuotaName}}", QuotaId:quotaId}, url);
    $("#add_quota_html").html(result);
    $("#add_post_html").modal("toggle")
}




/**
 * 删除资源配额方法
 * @param id
 * @return {*}
 */
function deleteQuota(id) {
    var url = "/api/quota/"+id;
    var result = del({}, url);
    result = JSON.stringify(result);
    return result
}


/**
 * 保存资源配额
 */
function saveQuota(quotaId,detail) {
    if(!quotaId){
        quotaId = 0
    }
    var data = get_form_data();
    data["QuotaId"] = parseInt(quotaId);
    if(!checkValue(data,"QuotaName,QuotaCpu,QuotaMemory,Description,AppNumber,JobNumber,RegistryGroupNumber,PipelineNumber,LbNumber,ServiceNumber,PodNumber")){
        return
    }
    var cpu = $("input[name='QuotaCpu']");
    if(data["QuotaCpu"] && !checkNumber(data, "QuotaCpu")){
        setInputError(cpu, "errmsg")
        return
    }else{
        setInputOk(cpu)
    }
    if(!data["UserName"] && !data["GroupName"]) {
        return
    }
    if(data["UserName"]) {
        data["UserName"] = data["UserName"].replace(/--请选择--/g,"");
    }
    if(data["GroupName"]) {
        data["GroupName"] = data["GroupName"].replace(/--请选择--/g,"");
    }
    var mem = $("input[name='QuotaMemory']")
    if(data["QuotaMemory"] && !checkNumber(data, "QuotaMemory")){
        setInputError(mem, "errmsg")
        return
    }else{
        setInputOk(mem)
    }
    var url = "/api/quota";
    var result = post(data, url);
    result = JSON.stringify(result);
    if (result.indexOf("保存成功") != -1){
        $("#add_post_html").modal("toggle")
        success(result)
        if (detail){
            window.location.reload();
        }else{
            toQuotaList();
        }
    }else{
        faild(result)
    }
}/**
 * Created by zhaoyun on 2018/1/5.
 */



function toQuotaList() {
    window.location.href = "/base/quota/list"
}

/**
 * 删除资源配额
 * 2018-02-12 06:55
 */
function deleteQuotaSwal(id) {
    Swal("删除资源配额", "warning", "确认操作", "不操作", "成功", "失败", " deleteQuota("+id+")", "topQuotaList()");
}
