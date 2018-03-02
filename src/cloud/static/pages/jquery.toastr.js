/**
 * Theme: Adminox Admin Template
 * Author: Coderthemes
 * Toastr js
 */

$("#toastr-one").click(function () {
    $.toast({
        heading: 'Heads up!',
        text: 'This alert needs your attention, but it is not super important.',
        position: 'top-right',
        loaderBg: '#3b98b5',
        icon: 'info',
        hideAfter: 3000,
        stack: 1
    });
});
$("#toastr-two").click(function () {
    $.toast({
        heading: 'Holy guacamole!',
        text: 'You should check in on some of those fields below.',
        position: 'top-right',
        loaderBg: '#da8609',
        icon: 'warning',
        hideAfter: 3000,
        stack: 1
    });
});
$("#toastr-three").click(function () {
    $.toast({
        heading: 'Well done!',
        text: 'You successfully read this important alert message.',
        position: 'top-right',
        loaderBg: '#5ba035',
        icon: 'success',
        hideAfter: 3000,
        stack: 1
    });
});
$("#toastr-four").click(function () {
    $.toast({
        heading: 'Oh snap!',
        text: 'Change a few things up and try submitting again.',
        position: 'top-right',
        loaderBg: '#bf441d',
        icon: 'error',
        hideAfter: 3000,
        stack: 1
    });
});

$("#toastr-five").click(function () {
    $.toast({
        heading: 'How to contribute?!',
        text: [
            'Fork the repository',
            'Improve/extend the functionality',
            'Create a pull request'
        ],
        position: 'top-right',
        loaderBg: '#1ea69a',
        hideAfter: 3000,
        stack: 1
    })
});

$("#toastr-six").click(function () {
    $.toast({
        heading: 'Can I add <em>icons</em>?',
        text: 'Yes! check this <a href="https://github.com/kamranahmedse/jquery-toast-plugin/commits/master">update</a>.',
        hideAfter: false,
        position: 'top-right',
        loaderBg: '#1ea69a',
        stack: 1
    })
});

$("#toastr-seven").click(function () {
    $.toast({
        text: 'Set the `hideAfter` property to false and the toast will become sticky.',
        hideAfter: false,
        position: 'top-right',
        loaderBg: '#1ea69a',
        stack: 1
    })
});

$("#toastr-eight").click(function () {
    $.toast({
        text: 'Set the `showHideTransition` property to fade|plain|slide to achieve different transitions',
        heading: 'Fade transition',
        showHideTransition: 'fade',
        position: 'top-right',
        loaderBg: '#1ea69a',
        hideAfter: 3000,
        stack: 1
    })
});

$("#toastr-nine").click(function () {
    $.toast({
        text: 'Set the `showHideTransition` property to fade|plain|slide to achieve different transitions',
        heading: 'Slide transition',
        showHideTransition: 'slide',
        position: 'top-right',
        loaderBg: '#1ea69a',
        hideAfter: 3000,
        stack: 1
    })
});

$("#toastr-ten").click(function () {
    $.toast({
        text: 'Set the `showHideTransition` property to fade|plain|slide to achieve different transitions',
        heading: 'Plain transition',
        showHideTransition: 'plain',
        position: 'top-right',
        loaderBg: '#1ea69a',
        hideAfter: 3000,
        stack: 1
    })
});
