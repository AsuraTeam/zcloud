
/**
* Theme: Adminox Admin Template
* Author: Coderthemes
* Ratings
*/

;(function ($) {
    $(function () {
        $('#default').raty({
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-warning'
        });

        $('#score').raty({
            score: 3,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#score-callback').raty({
            score: function () {
                return $(this).attr('data-score');
            }
        });

        $('#scoreName').raty({
            scoreName: 'entity[score]',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-warning'
        });

        $('#number').raty({
            number: 10,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#number-callback').raty({
            number: function () {
                return $(this).attr('data-number');
            }
        });

        $('#numberMax').raty({
            numberMax: 5,
            number: 100,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-purple'
        });

        $('#readOnly').raty({
            readOnly: true,
            score: 3,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-success'
        });

        $('#readOnly-callback').raty({
            readOnly: function () {
                return 'true becomes readOnly' == 'true becomes readOnly';
            }
        });

        $('#noRatedMsg').raty({
            readOnly: true,
            noRatedMsg: "I'am readOnly and I haven't rated yet!",
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#halfShow-true').raty({
            score: 3.26
        });

        $('#halfShow-false').raty({
            halfShow: false,
            score: 3.26,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#round').raty({
            round: {down: .26, full: .6, up: .76},
            score: 3.26,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-pink'
        });

        $('#half').raty({
            half: true
        });

        $('#starHalf').raty({
            half: true,
            starHalf: 'fa fa-star-half text-danger',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#click').raty({
            click: function (score, evt) {
                alert('ID: ' + $(this).attr('id') + "\nscore: " + score + "\nevent: " + evt.type);
            }
        });

        $('#hints').raty({hints: ['a', null, '', undefined, '*_*']});

        $('#star-off-and-star-on').raty({
            starOff: 'fa fa-bell-o text-muted',
            starOn: 'fa fa-bell text-custom'
        });

        $('#cancel').raty({
            cancel: true,
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#cancelHint').raty({
            cancel: true,
            cancelHint: 'My cancel hint!',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-success'
        });

        $('#cancelPlace').raty({
            cancel: true,
            cancelPlace: 'right',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-purple'
        });

        $('#cancel-off-and-cancel-on').raty({
            cancel: true,
            cancelOff: 'fa fa-minus-square-o text-muted',
            cancelOn: 'fa fa-minus-square text-danger'
        });

        $('#iconRange').raty({
            iconRange: [
                {range: 1, on: 'fa fa-cloud', off: 'fa fa-circle-o'},
                {range: 2, on: 'fa fa-cloud-download', off: 'fa fa-circle-o'},
                {range: 3, on: 'fa fa-cloud-upload', off: 'fa fa-circle-o'},
                {range: 4, on: 'fa fa-circle', off: 'fa fa-circle-o'},
                {range: 5, on: 'fa fa-cogs', off: 'fa fa-circle-o'}
            ]
        });

        $('#size-md').raty({
            cancel: true,
            half: true
        });

        $('#size-lg').raty({
            cancel: true,
            half: true
        });

        $('#target-div').raty({
            cancel: true,
            target: '#target-div-hint',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-custom'
        });

        $('#targetType').raty({
            cancel: true,
            target: '#targetType-hint',
            targetType: 'score',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-warning'
        });

        $('#targetFormat').raty({
            target: '#targetFormat-hint',
            targetFormat: 'Rating: {score}',
            starOff: 'fa fa-star-o text-muted',
            starOn: 'fa fa-star text-danger'
        });

        $('#mouseover').raty({
            mouseover: function (score, evt) {
                alert('ID: ' + $(this).attr('id') + "\nscore: " + score + "\nevent: " + evt.type);
            }
        });

        $('#mouseout').raty({
            width: 150,
            mouseout: function (score, evt) {
                alert('ID: ' + $(this).attr('id') + "\nscore: " + score + "\nevent: " + evt.type);
            }
        });
    });
})(jQuery);