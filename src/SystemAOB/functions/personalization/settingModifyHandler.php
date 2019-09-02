<?php
include_once("../../../auth.php");
include_once("../user/userIsolation.php");
$configName = $_POST["autoConfigBaseConfigurationFilename"];
$configFile = "sysconf/" . $configName . ".config";
if (file_exists($configFile)){
	
}else if (file_exists($userConfigDirectory . "SystemAOB/sysconf/" . $configName . ".config")){
	$configFile = $userConfigDirectory . "SystemAOB/sysconf/" . $configName . ".config";
}else{
	die("ERROR. Required config file not found in default config setting directory.");
}

$config = file_get_contents($configFile);
$config = json_decode($config,true);
//As checkbox will not pass through values, default values for all boolean should be false.
foreach ($config as $key => $value) {
	if ($config[$key][2] == "boolean"){
		$config[$key][3] = "false";
	}
}
foreach ($_POST as $key => $value) {
	if ($key != "autoConfigBaseConfigurationFilename"){
		//This is not the system defined config value. Write the rest of the content in.
		if ($config[$key][2] == "boolean"){
			//This is boolean. Change this from "on" to "true" if this value received
			$config[$key][3] = "true";
		}else{
			$config[$key][3] = $value;
		}
	}
}

file_put_contents($configFile,json_encode($config));
header("Location: autoConfig.php?configName=" . $configName . "&update=" . time());
exit(0);

?>