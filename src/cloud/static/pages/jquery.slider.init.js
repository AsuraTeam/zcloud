/**
 * Theme: Adminox Dashboard
 * Author: Coderthemes
 * Slider
 */

$('.single-item').slick({
    autoplay: true,
    autoplaySpeed: 2000
});

$('.multiple-items').slick({
    infinite: true,
    slidesToShow: 3,
    slidesToScroll: 1
});

$('.responsive-slider').slick({
    dots: true,
    infinite: false,
    speed: 300,
    slidesToShow: 4,
    slidesToScroll: 4,
    responsive: [
        {
            breakpoint: 1024,
            settings: {
                slidesToShow: 3,
                slidesToScroll: 3,
                infinite: true,
                dots: true
            }
        },
        {
            breakpoint: 600,
            settings: {
                slidesToShow: 2,
                slidesToScroll: 2
            }
        },
        {
            breakpoint: 480,
            settings: {
                slidesToShow: 1,
                slidesToScroll: 1
            }
        }
        // You can unslick at a given breakpoint now by adding:
        // settings: "unslick"
        // instead of a settings object
    ]
});

$('.center-item-slider').slick({
    centerMode: true,
    centerPadding: '60px',
    slidesToShow: 3,
    responsive: [
        {
            breakpoint: 768,
            settings: {
                arrows: false,
                centerMode: true,
                centerPadding: '40px',
                slidesToShow: 3
            }
        },
        {
            breakpoint: 480,
            settings: {
                arrows: false,
                centerMode: true,
                centerPadding: '40px',
                slidesToShow: 1
            }
        }
    ]
});

$('.lazy-load-slider').slick({
    lazyLoad: 'ondemand',
    slidesToShow: 3,
    slidesToScroll: 1
});

$('.fade-slider').slick({
    dots: true,
    infinite: true,
    speed: 500,
    fade: true,
    cssEase: 'linear',
    autoplay: true,
    autoplaySpeed: 2000
});


// Slider syncing
$('.slider-syncing-for').slick({
    slidesToShow: 1,
    slidesToScroll: 1,
    arrows: false,
    fade: true,
    asNavFor: '.slider-syncing-nav'
});
$('.slider-syncing-nav').slick({
    slidesToShow: 5,
    slidesToScroll: 1,
    asNavFor: '.slider-syncing-for',
    dots: true,
    centerMode: true,
    focusOnSelect: true
});