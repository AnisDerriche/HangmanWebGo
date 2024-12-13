var button = document.getElementById('choix');

button.addEventListener('mouseover', function() {
    button.src = '../static/Image/play_button400.png';
});

button.addEventListener('mouseout', function() {
    button.src = '../static/Image/play_button2_400.png';
});

document.querySelectorAll('.jumpable').forEach(function(img) {
    img.addEventListener('mouseover', function() {
        img.classList.add('jump');
        setTimeout(function() {
            img.classList.remove('jump');
        }, 2000); // Retire la classe après 2 secondes
    });
});
