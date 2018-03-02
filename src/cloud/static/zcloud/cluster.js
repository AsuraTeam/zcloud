

!function ($) {
    "use strict";

    var SweetAlert = function () {
    };

    //examples
    SweetAlert.prototype.init = function () {
        //Parameter
        $('#deleteClusterId').click(function () {
            swal({
                title: '删除该集群',
                text: "",
                type: 'warning',
                showCancelButton: true,
                confirmButtonText: '确认删除',
                cancelButtonText: '不删除',
                confirmButtonClass: 'btn btn-success',
                cancelButtonClass: 'btn btn-danger m-l-10',
                buttonsStyling: false
            }).then(function () {
                var result = deleteCluster()
                if (result.indexOf("删除成功") != -1){
                    // success(result)
                    swal(
                        '删除成功!',
                        result,
                        'success'
                    )
                    window.location.href = "/base/cluster/list"
                }else{
                    swal(
                        '删除失败!',
                        result,
                        'error'
                    )
                }

            }, function (dismiss) {
                // dismiss can be 'cancel', 'overlay',
                // 'close', and 'timer'
            })
        });
    },
        $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
}(window.jQuery),

//initializing
    function ($) {
        "use strict";
        $.SweetAlert.init()
    }(window.jQuery);

/**
 * Created by zhaoyun on 2018/1/3.
 */
