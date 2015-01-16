# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def threadFunction_1():
	global i
	for j in range(0,1000000):	
		i+=1
def threadFunction_2():
	global i
	for j in range(0,1000000):
		i-=1

# Potentially useful thing:
# In Python you "import" a global variable, instead of "export"ing it when you declare it
# (This is probably an effort to make you feel bad about typing the word "global")


def main():
	thread_1 = Thread(target = threadFunction_1, args = (),)
	thread_1.start()
	thread_2 = Thread(target = threadFunction_2, args = (),)
	thread_2.start()

	thread_1.join()
	thread_2.join()
	print(i)

main()
