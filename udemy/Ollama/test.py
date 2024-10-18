def Prime(x):
    a = True
    b = 1
    while (b < x)and(a == True):
        if not(x%b==0 and x!=2*b):
            b += 1
        else:
            print(str(x)+" is prime")
            break
    else: 
        print("This number "+ str(x)+ " is not a prime number. Try again!")

Prime(int(input("Insert a positive number to check: ")))
