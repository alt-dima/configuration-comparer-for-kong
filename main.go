package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"slices"

	"github.com/kong/go-kong/kong"
)

func slicePointersToValues(slicePointers []*string) []string {
	sliceValues := make([]string, len(slicePointers))
	for i, s := range slicePointers {
		sliceValues[i] = *s
	}
	slices.Sort(sliceValues)
	return sliceValues
}

func main() {
	clientUrl1 := &os.Args[1]
	clientUrl2 := &os.Args[2]

	client1, err := kong.NewClient(clientUrl1, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client2, err := kong.NewClient(clientUrl2, nil)
	if err != nil {
		log.Fatalln(err)
	}

	//compare routes and services
	allRoutesClient1, err := client1.Routes.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}
	allRoutesClient2, err := client2.Routes.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Amount routes in %v: %v \n", *clientUrl1, len(allRoutesClient1))
	fmt.Printf("Amount routes in %v: %v \n", *clientUrl2, len(allRoutesClient2))

	for _, client1Route := range allRoutesClient1 {
		result := false

		client1RoutePaths := slicePointersToValues(client1Route.Paths)
		strClient1RoutePaths := fmt.Sprint(client1RoutePaths)

		for _, client2Route := range allRoutesClient2 {

			if slices.Equal(client1RoutePaths, slicePointersToValues(client2Route.Paths)) {
				result = true
				//checking route parameters
				if *client1Route.PreserveHost != *client2Route.PreserveHost {
					fmt.Printf("Route %v PreserveHost %v not equals %v \n", strClient1RoutePaths, *client1Route.PreserveHost, *client2Route.PreserveHost)
				}
				if *client1Route.StripPath != *client2Route.StripPath {
					fmt.Printf("Route %v StripPath %v not equals %v \n", strClient1RoutePaths, *client1Route.StripPath, *client2Route.StripPath)
				}
				if !slices.Equal(slicePointersToValues(client1Route.Methods), slicePointersToValues(client2Route.Methods)) {
					fmt.Printf("Route %v Methods not equals \n", strClient1RoutePaths)
				}
				if !slices.Equal(slicePointersToValues(client1Route.Hosts), slicePointersToValues(client2Route.Hosts)) {
					fmt.Printf("Route %v Hosts not equals \n", strClient1RoutePaths)
				}
				if !slices.Equal(slicePointersToValues(client1Route.Protocols), slicePointersToValues(client2Route.Protocols)) {
					fmt.Printf("Route %v Protocols not equals \n", strClient1RoutePaths)
				}

				//check route plugins
				allPluginsRouteClient1, err := client1.Plugins.ListAllForRoute(nil, client1Route.ID)
				if err != nil {
					log.Fatalln(err)
				}
				allPluginsRouteClient2, err := client2.Plugins.ListAllForRoute(nil, client2Route.ID)
				if err != nil {
					log.Fatalln(err)
				}
				for _, pluginRouteClient1 := range allPluginsRouteClient1 {
					result := false

					for _, pluginRouteClient2 := range allPluginsRouteClient2 {
						if *pluginRouteClient1.Name == *pluginRouteClient2.Name && *pluginRouteClient1.Enabled == *pluginRouteClient2.Enabled {
							result = true
							delete(pluginRouteClient1.Config, "anonymous")
							delete(pluginRouteClient2.Config, "anonymous")
							delete(pluginRouteClient1.Config, "okta_consumer")
							delete(pluginRouteClient2.Config, "okta_consumer")
							if !reflect.DeepEqual(pluginRouteClient1.Config, pluginRouteClient2.Config) {
								//fmt.Println(pluginRouteClient1.Config)
								//fmt.Println(pluginRouteClient2.Config)
								fmt.Println("Route " + strClient1RoutePaths + " plugin config " + *pluginRouteClient1.Name + " not equals in " + *clientUrl2)
							}
							break
						}
					}
					if result == false {
						fmt.Println("Route " + strClient1RoutePaths + " plugin " + *pluginRouteClient1.Name + " does not exists in " + *clientUrl2)
					}
				}

				//check route service
				serviceRouteClient1, err := client1.Services.GetForRoute(nil, client1Route.ID)
				if err != nil {
					log.Fatalln(err)
				}
				serviceRouteClient2, err := client2.Services.GetForRoute(nil, client2Route.ID)
				if err != nil {
					log.Fatalln(err)
				}
				if *serviceRouteClient1.Host != *serviceRouteClient2.Host {
					fmt.Printf("Route %v target service Host not equals \n", strClient1RoutePaths)
				}
				if *serviceRouteClient1.Port != *serviceRouteClient2.Port {
					fmt.Printf("Route %v target service Port not equals \n", strClient1RoutePaths)
				}
				if *serviceRouteClient1.Path != *serviceRouteClient2.Path {
					fmt.Printf("Route %v target service Path not equals \n", strClient1RoutePaths)
				}

				//check service plugins
				allPluginsServiceClient1, err := client1.Plugins.ListAllForService(nil, serviceRouteClient1.ID)
				if err != nil {
					log.Fatalln(err)
				}
				allPluginsServiceClient2, err := client2.Plugins.ListAllForService(nil, serviceRouteClient2.ID)
				if err != nil {
					log.Fatalln(err)
				}
				for _, pluginServiceClient1 := range allPluginsServiceClient1 {
					result := false

					for _, pluginServiceClient2 := range allPluginsServiceClient2 {
						if *pluginServiceClient1.Name == *pluginServiceClient2.Name && *pluginServiceClient1.Enabled == *pluginServiceClient2.Enabled {
							result = true
							delete(pluginServiceClient1.Config, "append")
							delete(pluginServiceClient2.Config, "append")
							//delete(pluginServiceClient1.Config, "okta_consumer")
							//delete(pluginServiceClient2.Config, "okta_consumer")
							if !reflect.DeepEqual(pluginServiceClient1.Config, pluginServiceClient2.Config) {
								//fmt.Println(pluginServiceClient1.Config)
								//fmt.Println(pluginServiceClient2.Config)
								fmt.Println("Service " + *serviceRouteClient1.Name + " plugin config " + *pluginServiceClient1.Name + " not equals in " + *clientUrl2)
							}
							break
						}
					}
					if result == false {
						fmt.Println("Route " + strClient1RoutePaths + " plugin " + *pluginServiceClient1.Name + " does not exists in " + *clientUrl2)
					}
				}

				break
			}
		}
		if result == false {
			fmt.Println("Route " + strClient1RoutePaths + " does not exists in " + *clientUrl2)
		}
	}

	//compare consumers
	allConsumersClient1, err := client1.Consumers.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}
	allConsumersClient2, err := client2.Consumers.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Amount consumers in %v: %v \n", *clientUrl1, len(allConsumersClient1))
	fmt.Printf("Amount consumers in %v: %v \n", *clientUrl2, len(allConsumersClient2))

	fmt.Printf("Only Consumers plugins and ACLs comparison currently implemented!!! \n")

	for _, client1Consumer := range allConsumersClient1 {
		result := false

		for _, client2Consumer := range allConsumersClient2 {

			if *client1Consumer.Username == *client2Consumer.Username && *client1Consumer.CustomID == *client2Consumer.CustomID {
				result = true

				//check consumer plugins
				allPluginsConsumerClient1, err := client1.Plugins.ListAllForConsumer(nil, client1Consumer.ID)
				if err != nil {
					log.Fatalln(err)
				}
				allPluginsConsumerClient2, err := client2.Plugins.ListAllForConsumer(nil, client2Consumer.ID)
				if err != nil {
					log.Fatalln(err)
				}
				for _, pluginConsumerClient1 := range allPluginsConsumerClient1 {
					result := false

					for _, pluginConsumerClient2 := range allPluginsConsumerClient2 {
						if *pluginConsumerClient1.Name == *pluginConsumerClient2.Name && *pluginConsumerClient1.Enabled == *pluginConsumerClient2.Enabled {
							result = true
							//fmt.Println(pluginConsumerClient1.Config)
							//delete(pluginConsumerClient1.Config, "append")
							//delete(pluginConsumerClient2.Config, "append")
							//delete(pluginConsumerClient1.Config, "okta_consumer")
							//delete(pluginConsumerClient2.Config, "okta_consumer")
							if !reflect.DeepEqual(pluginConsumerClient1.Config, pluginConsumerClient2.Config) {
								//fmt.Println(pluginConsumerClient1.Config)
								//fmt.Println(pluginConsumerClient2.Config)
								fmt.Println("Consumer " + *client1Consumer.Username + " plugin config " + *pluginConsumerClient1.Name + " not equals in " + *clientUrl2)
							}
							break
						}
					}
					if result == false {
						fmt.Println("Consumer " + *client1Consumer.Username + " plugin " + *pluginConsumerClient1.Name + " does not exists in " + *clientUrl2)
					}
				}

				//check consumer ACLs
				allACLsConsumerClient1, _, err := client1.ACLs.ListForConsumer(nil, client1Consumer.ID, nil)
				if err != nil {
					log.Fatalln(err)
				}
				allACLsConsumerClient2, _, err := client2.ACLs.ListForConsumer(nil, client2Consumer.ID, nil)
				if err != nil {
					log.Fatalln(err)
				}
				for _, aclConsumerClient1 := range allACLsConsumerClient1 {
					result := false

					for _, aclConsumerClient2 := range allACLsConsumerClient2 {
						if *aclConsumerClient1.Group == *aclConsumerClient2.Group {
							result = true

							break
						}
					}
					if result == false {
						fmt.Println("Consumer " + *client1Consumer.Username + " ACL " + *aclConsumerClient1.Group + " does not exists in " + *clientUrl2)
					}
				}
				break
			}
		}
		if result == false {
			fmt.Println("Consumer " + *client1Consumer.Username + " does not exists in " + *clientUrl2)
		}
	}

	//compare global plugins
	allPluginsClient1, err := client1.Plugins.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}

	var allGlobalPluginsClient1 []*kong.Plugin

	for _, PluginClient1 := range allPluginsClient1 {
		if PluginClient1.Route == nil && PluginClient1.Service == nil && PluginClient1.Consumer == nil {
			allGlobalPluginsClient1 = append(allGlobalPluginsClient1, PluginClient1)
		}
	}

	allPluginsClient2, err := client2.Plugins.ListAll(nil)
	if err != nil {
		log.Fatalln(err)
	}

	var allGlobalPluginsClient2 []*kong.Plugin

	for _, PluginClient2 := range allPluginsClient2 {
		if PluginClient2.Route == nil && PluginClient2.Service == nil && PluginClient2.Consumer == nil {
			allGlobalPluginsClient2 = append(allGlobalPluginsClient2, PluginClient2)
		}
	}

	fmt.Printf("Amount global plugins in %v: %v \n", *clientUrl1, len(allGlobalPluginsClient1))
	fmt.Printf("Amount global plugins in %v: %v \n", *clientUrl2, len(allGlobalPluginsClient2))

	for _, globalPluginClient1 := range allGlobalPluginsClient1 {
		result := false

		for _, globalPluginClient2 := range allGlobalPluginsClient2 {
			if *globalPluginClient1.Name == *globalPluginClient2.Name && *globalPluginClient1.Enabled == *globalPluginClient2.Enabled {
				result = true
				delete(globalPluginClient1.Config, "anonymous")
				delete(globalPluginClient2.Config, "anonymous")
				delete(globalPluginClient1.Config, "okta_consumer")
				delete(globalPluginClient2.Config, "okta_consumer")
				if !reflect.DeepEqual(globalPluginClient1.Config, globalPluginClient2.Config) {
					//fmt.Println(pluginRouteClient1.Config)
					//fmt.Println(pluginRouteClient2.Config)
					fmt.Println("Global plugin config " + *globalPluginClient1.Name + " not equals in " + *clientUrl2)
				}
				break
			}
		}
		if result == false {
			fmt.Println("Global plugin " + *globalPluginClient1.Name + " does not exists in " + *clientUrl2)
		}
	}

	// //compare consumer groups
	// allConsGroupsClient1, err := client1.ConsumerGroups.ListAll(nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// allConsGroupsClient2, err := client2.ConsumerGroups.ListAll(nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Printf("Amount consumer groups in %v: %v \n", *clientUrl1, len(allConsGroupsClient1))
	// fmt.Printf("Amount consumer groups in %v: %v \n", *clientUrl2, len(allConsGroupsClient2))

	// for _, consGroupClient1 := range allConsGroupsClient1 {
	// 	result := false

	// 	for _, consGroupClient2 := range allConsGroupsClient2 {
	// 		if *consGroupClient1.Name == *consGroupClient2.Name && *consGroupClient1.ID == *consGroupClient1.ID {
	// 			result = true

	// 			//check consumer plugins
	// 			consumersGroupClient1, err := client1.ConsumerGroupConsumers.ListAll(nil, consGroupClient1.ID)
	// 			if err != nil {
	// 				log.Fatalln(err)
	// 			}
	// 			consumersGroupClient2, err := client2.ConsumerGroupConsumers.ListAll(nil, consGroupClient2.ID)
	// 			if err != nil {
	// 				log.Fatalln(err)
	// 			}
	// 			for _, groupConsumerClient1 := range consumersGroupClient1.Consumers {
	// 				result := false

	// 				for _, groupConsumerClient2 := range consumersGroupClient2.Consumers {
	// 					if *groupConsumerClient1.Username == *groupConsumerClient2.Username && *groupConsumerClient1.ID == *groupConsumerClient2.ID {
	// 						result = true

	// 						break
	// 					}
	// 				}
	// 				if result == false {
	// 					fmt.Println("Consumer " + *consGroupClient1.Name + " plugin " + *groupConsumerClient1.Username + " does not exists in " + *clientUrl2)
	// 				}
	// 			}
	// 			break
	// 		}
	// 	}
	// 	if result == false {
	// 		fmt.Println("Cons Group " + *consGroupClient1.Name + " does not exists in " + *clientUrl2)
	// 	}
	// }
}
