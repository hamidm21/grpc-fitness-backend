# Workflow for authentication 
  
Step one:  
Request activation code from PhoneSignIn api using phone number.  
Step two:  
Request for a jwt token from Login api using phone number & activation code.  
Step three:  
Add the obtained jwt token from previous step to every api call metadata using header `x-token`.