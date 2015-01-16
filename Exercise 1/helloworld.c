// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

// Note the return type: void*

int i = 0;

void* threadFunction_1(){
	for (int j=0;j<1000000;j++){
		i++;
	}
	//printf("%d",i);
	return NULL;
}

void* threadFunction_2(){
	for (int k=0;k<1000000;k++){
		i--;
	}
	//printf("%d",i);
	return NULL;
}

int main(){
	
	
	pthread_t thread_1;
	pthread_t thread_2;
	
	pthread_create(&thread_1, NULL, threadFunction_1, NULL);
	pthread_create(&thread_2, NULL, threadFunction_2, NULL);
	// Arguments to a thread would be passed here ---------^

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);
	printf("%d", i);
	return 0;
}
