<?php


$url = "http://127.0.0.1:8088/validate/";

$ip = getenv('HTTP_CLIENT_IP')?:
getenv('HTTP_X_FORWARDED_FOR')?:
getenv('HTTP_X_FORWARDED')?:
getenv('HTTP_FORWARDED_FOR')?:
getenv('HTTP_FORWARDED')?:
getenv('REMOTE_ADDR');

echo $ip."<br>";

$headers = get_headers( $url.$ip, 1 );


if ($headers[0] == "HTTP/1.1 204 No Content"){
    echo "Your IP is Not blocked. Show ADs";
    // тут что угодно
    echo '<script type="text/javascript">
    alert("Hello! I am an alert box!!");
    </script>';
}

if ($headers[0] == "HTTP/1.1 200 OK"){
    echo "Your IP is Blocked. Forbidden Show ADs";
    echo '<script type="text/javascript">
          alert("Hello! I am an alert box!!");
    </script>';
}
