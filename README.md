<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  
</head>
<body>

  <h1 align="center">Train-Services</h1>

  <p align="center">
    <a href="https://github.com/lucilu1320/Train-Services"><img src="https://img.shields.io/github/stars/lucilu1320/Train-Services?style=flat-square" alt="Stars"></a>
    <a href="https://github.com/lucilu1320/Train-Services/issues"><img src="https://img.shields.io/github/issues/lucilu1320/Train-Services?style=flat-square" alt="Issues"></a>
    <a href="https://github.com/lucilu1320/Train-Services/blob/master/LICENSE"><img src="https://img.shields.io/github/license/lucilu1320/Train-Services?style=flat-square" alt="License"></a>
  </p>

  <h2>Getting Started</h2>

  <h3>Prerequisites:</h3>

  <ul>
    <li>Golang must be installed</li>
    <li>Protocol Buffer must be installed</li>
  </ul>



  <h3>Installation:</h3>

  <ul>
    <li>To install Protocol Buffer on both Linux and macOS, follow the instructions provided in the link/reference:(both for linux and MacOS): https://grpc.io/docs/protoc-installation/</li>
    <li>Please note that in case of MacOS, Homebrew must be installed on your Mac system.</li>
  </ul>

  <h4>To install Homebrew:(if needed)</h4>
  <p> Install Command-line Tools for Xcode:
    Open the Terminal app (located in /Applications/Utilities).</p>

  <ol>
    
  <li>Run the following command to install the command-line tools:
      <code>xcode-select --install</code> </li>
    <li>Install Homebrew:
      <code>/bin/bash -c "<span class="math-inline">\(curl \-fsSL https\://raw\.githubusercontent\.com/Homebrew/install/HEAD/install\.sh\)"</code> </li>
        
<li>Verify Installation:
<code>brew --version</code>
</li>
<li>Set PATH Environment Variable (Apple Silicon Macs only):
<code> echo 'eval "</span>(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile </code>
      <code>source ~/.zprofile </code>
    </li>
<li> Check the URL below for installation of Protocol Buffer now using brew: https://grpc.io/docs/protoc-installation/</li>
  </ol>

<h4>Choco Installation for windows:</h4>
<ol>
<li> If you want to install on windows, install using choco following the link: https://chocolatey.org/install </li>
<li>Now use the command below to install Protocol Buffer on windows using choco: <code> choco install protoc </code> </li>

  
</ol>

  <h2>Usage:</h2>

  <ol>
    <li>Clone the repository:
      <code>git clone https://github.com/lucilu1320/Train-Services</code>
    </li>
    <li>Navigate to the project directory:
      <code>cd Train-Services</code>
      
   </li>
   <li>Run the command to install the modules:
     <code>go mod tidy</code>
   </li>
    <li>Open two terminals:
      <ul>
        <li>In terminal 1, navigate to the client directory:
          <code>cd client</code>
        </li>
        <li>In terminal 2, navigate to the server directory:
          <code>cd server</code>
        </li>
      </ul>
    </li>
    <li>Run the main.go files in each terminal:
      <ul>
        <li>Server: <code>go run *.go</code></li>
        <li>Client: <code>go run *.go</code></li>
        <li>Enjoy!</li>
        
  </ul>
    </li>
  </ol>
<h2>----------------------------------END--------------------------------------</h2>
