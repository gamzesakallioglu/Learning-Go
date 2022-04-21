    
    
   ### Book Archive - Attainments

- Arguments           
We use command-line arguments to parametize execution of programs. In order to access the arguments, type **os.Args**. This is a slice that the first element is the path to the program. (temp file for run command) **os.Args[1:]** is the way to access to given arguments since the first element won't be included.        
         
- Iterate a slice              
It's easy to iterate through a slice with **for** and **range**         
    for i, v := range(slice1){       
       
    }      
i is for index and v is for value    
