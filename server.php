<?php

print_r($_REQUEST);
echo PHP_EOL;
print_r($_FILES);
echo PHP_EOL;
foreach ($_FILES as $info) {
    move_uploaded_file($info['tmp_name'], 'uploads/' . $info['name']);
}
if (isset($_GET['download']) && $_GET['download'] !== '') {
    $file = $_GET['download'];
    header('Content-Description: File Transfer');
    header('Content-Type: application/octet-stream');
    header('Content-Disposition: attachment; filename=' . basename($file));
    header('Content-Transfer-Encoding: binary');
    header('Expires: 0');
    header('Cache-Control: must-revalidate, post-check=0, pre-check=0');
    header('Pragma: public');
    header('Content-Length: ' . filesize($file));
    ob_clean();
    flush();
    readfile($file);
}