<?php

$redis_handle = new Redis();
$redis_handle->connect('127.0.0.1', 8299, 10);
$result = $redis_handle->rawCommand('filter', '妈的-我看他说话的语气，好屌啊');
echo $result . PHP_EOL;
$result = $redis_handle->rawCommand('add', '好屌');
var_dump($result);
$result = $redis_handle->rawCommand('add', '好屌');
var_dump($result);
$result = $redis_handle->rawCommand('filter', '妈的-我看他说话的语气，好屌啊');
echo $result . PHP_EOL;