# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread
import threading

i = 0
lock = threading.Lock()

def threadFunction_1():
	global i
	lock.acquire()	
	for j in range(0,1000001):	
		i+=1
	lock.release()

def threadFunction_2():
	global i
	lock.acquire()
	for j in range(0,1000000):
		i-=1
	lock.release()

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
