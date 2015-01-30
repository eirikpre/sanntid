// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

// Note the return type: void*

int i = 0;
pthread_mutex_t lock;

void* threadFunction_1(){
	pthread_mutex_lock(&lock);

	for (int j=0;j<1000000;j++){
		i++;
	}
	//printf("%d",i);
	pthread_mutex_unlock(&lock);
	return NULL;
}

void* threadFunction_2(){
	pthread_mutex_lock(&lock);
		
	for (int k=0;k<100000;k++){
		i--;
	}	
	//printf("%d",i);
	pthread_mutex_unlock(&lock);	
	return NULL;
}

int main(){
	
	pthread_mutex_init(&lock,NULL);
	
	pthread_t thread_1;
	pthread_t thread_2;
	pthread_create(&thread_1, NULL, threadFunction_1, NULL);
	pthread_create(&thread_2, NULL, threadFunction_2, NULL);

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);
	pthread_mutex_destroy(&lock);

	printf("%d \n", i);

	return 0;
}
