# Modem Functionality

## Modem Test Procedure
There are two challenge procedures, then an error checking procedure.

### Challenge Procedures
Challenge procedure 1 requires the following bytes be sent in sequence:
* 0x32
* 0x3C
* 0x46
* At this point, the Modem driver will send identity data to the server. Subsequently the server must respond with
  * 0x27 - means the identity data was ok
  * 0x63 - means the identity data was bad

Challenge procedure 2 will try to read 1 byte five times before it will error out. This procedure will accept two types of input:
* 0x7 - this will report success and move to next step
* 0x21 - this will cause the procedure to loop without erroring

### Error Checking Procedures
Next comes the error checking procedures.
First, Modem driver will send all bytes between 0x00 and 0xFF inclusive.

Server must respond with the following byte:
* 0x04

Anything else will result in a "Bad Response" error

Next, Modem driver expects server to send all bytes between 0x00 and 0xFF inclusive. The Modem driver allows for up to 3 errors before it reports "Too many errors" and aborts the test.

Otherwise, the Modem driver will report "Test successful".

Finally, the Modem driver expects server to send an updated "phone number" which it will save for the next connection.
