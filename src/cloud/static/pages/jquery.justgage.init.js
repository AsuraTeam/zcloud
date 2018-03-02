/**
 * Theme: Adminox Admin Template
 * Author: Coderthemes
 * Module/App: Justgage
 */



document.addEventListener("DOMContentLoaded", function(event) {

    var g1, g2, g3, g4, g5, g6, g7, g8;

    var g1 = new JustGage({
        id: "g1",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Custom Width",
        label: "miles traveled",
        gaugeWidthScale: 0.2
    });

    var g2 = new JustGage({
        id: "g2",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Custom Shadow",
        label: "",
        shadowOpacity: 1,
        shadowSize: 10,
        shadowVerticalOffset: 5
    });

    var g3 = new JustGage({
        id: "g3",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Custom Colors",
        label: "",
        levelColors: [
            "#00fff6",
            "#ff00fc",
            "#1200ff"
        ]
    });

    var g4 = new JustGage({
        id: "g4",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Hide Labels",
        hideMinMax: true
    });


    var g5 = new JustGage({
        id: "g5",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Animation Type",
        label: "",
        startAnimationTime: 2000,
        startAnimationType: ">",
        refreshAnimationTime: 1000,
        refreshAnimationType: "bounce"
    });

    var g6 = new JustGage({
        id: "g6",
        value: getRandomInt(0, 100),
        min: 0,
        max: 100,
        title: "Minimal",
        label: "",
        hideMinMax: true,
        gaugeColor: "#fff",
        levelColors: ["#000"],
        hideInnerShadow: true,
        startAnimationTime: 1,
        startAnimationType: "linear",
        refreshAnimationTime: 1,
        refreshAnimationType: "linear"
    });

    var g7 = new JustGage({
        id: "g7",
        value: 72,
        min: 0,
        max: 100,
        donut: true,
        gaugeWidthScale: 0.6,
        counter: true,
        hideInnerShadow: true
    });

    var g8 = new JustGage({
        id: "g8",
        value : 72.15,
        min: 0,
        max: 100,
        decimals: 2,
        gaugeWidthScale: 0.6,
        customSectors: [{
            color : "#00ff00",
            lo : 0,
            hi : 50
        },{
            color : "#ff0000",
            lo : 50,
            hi : 100
        }],
        counter: true
    });

    document.getElementById('g8_refresh').addEventListener('click', function() {
        g8.refresh(getRandomInt(0, 100));
    });

    setInterval(function() {
        g1.refresh(getRandomInt(0, 100));
        g2.refresh(getRandomInt(0, 100));
        g3.refresh(getRandomInt(0, 100));
        g4.refresh(getRandomInt(0, 100));
        g5.refresh(getRandomInt(0, 100));
        g6.refresh(getRandomInt(0, 100));
        g7.refresh(getRandomInt(0, 100));
    }, 2500);

});
