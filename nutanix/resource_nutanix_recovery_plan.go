package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixRecoveryPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixRecoveryPlanCreate,
		Read:   resourceNutanixRecoveryPlanRead,
		Update: resourceNutanixRecoveryPlanUpdate,
		Delete: resourceNutanixRecoveryPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stage_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stage_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"delay_time_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"stage_work": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recover_entities": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"entity_info_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"any_entity_reference_kind": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"any_entity_reference_uuid": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"any_entity_reference_name": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"categories": categoriesSchema(),
															"script_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_script_exec": {
																			Type:     schema.TypeBool,
																			Required: true,
																		},
																		"timeout": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"parameters": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"floating_ip_assignment_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"vm_ip_assignment_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"test_floating_ip_config": {
													Type:     schema.TypeList,
													MinItems: 1,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"should_allocate_dynamically": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"recovery_floating_ip_config": {
													Type:     schema.TypeList,
													MinItems: 1,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"should_allocate_dynamically": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"vm_reference": {
													Type:     schema.TypeList,
													MinItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": {
																Type:     schema.TypeString,
																Required: true,
															},
															"uuid": {
																Type:     schema.TypeString,
																Required: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"vm_nic_information": {
													Type:     schema.TypeList,
													MinItems: 1,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"uuid": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"network_mapping_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_network_mapping_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"availability_zone_url": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"recovery_network": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MinItems: 1,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"virtual_network_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"vpc_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"subnet_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"gateway_ip": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"external_connectivity_state": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																	},
																},
															},
															"use_vpc_reference": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
												"test_network": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MinItems: 1,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"virtual_network_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"vpc_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"subnet_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"gateway_ip": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"external_connectivity_state": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Required: true,
																		},
																	},
																},
															},
															"use_vpc_reference": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixRecoveryPlanCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.RecoveryPlanInput{}
	spec := &v3.RecoveryPlanSpec{}
	metadata := &v3.Metadata{}
	recoveryPlan := &v3.RecoveryPlanResources{}

	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	if !nok {
		return fmt.Errorf("please provide the required attributes `name`")
	}

	if err := getMetadataAttributes(d, metadata, "recovery_plan"); err != nil {
		return err
	}

	getRecoveryPlanResources(d, recoveryPlan)

	if descok {
		spec.Description = desc.(string)
	}

	recoveryUUID, err := resourceNutanixRecoveryPlanExists(conn, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error checking if recovery_plan already exists %+v", err)
	}

	if recoveryUUID != nil {
		return fmt.Errorf("recovery_plan already with name %s exists , UUID %s", d.Get("name").(string), *recoveryUUID)
	}

	spec.Name = n.(string)
	spec.Resources = recoveryPlan
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateRecoveryPlan(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix RecoveryPlan %s: %+v", spec.Name, err)
	}

	d.SetId(*resp.Metadata.UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the RecoveryPlans to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for recovery_plan id (%s) to create: %+v", id, err)
	}

	// Setting Description because in Get request is not present.
	d.Set("description", resp.Spec.Description)

	return resourceNutanixRecoveryPlanRead(d, meta)
}

func resourceNutanixRecoveryPlanRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	resp, err := conn.V3.GetRecoveryPlan(id)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValuesList(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValuesList(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return err
	}
	if err := d.Set("stage_list", flattenStageList(resp.Spec.Resources.StageList)); err != nil {
		return err
	}
	if err := d.Set("parameters", flattenParameters(resp.Spec.Resources.Parameters)); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixRecoveryPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.RecoveryPlanInput{}
	metadata := &v3.Metadata{}
	res := &v3.RecoveryPlanResources{}
	spec := &v3.RecoveryPlanSpec{}

	id := d.Id()
	response, err := conn.V3.GetRecoveryPlan(id)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "RECOVERY_PLAN_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for protection rule id (%s) :%+v", id, err)
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec
		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").([]interface{})
		metadata.OwnerReference = validateRefList(or, utils.StringPtr("recovery_plan"))
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").([]interface{})
		metadata.ProjectReference = validateRefList(pr, utils.StringPtr("recovery_plan"))
	}
	if d.HasChange("name") {
		spec.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		spec.Description = d.Get("description").(string)
	}
	if d.HasChange("stage_list") {
		spec.Resources.StageList = expandStageList(d)
	}
	if d.HasChange("parameters") {
		spec.Resources.Parameters = expandParameters(d)
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateRecoveryPlan(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error recovery_plan subnet id %s): %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for protection rule (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixRecoveryPlanRead(d, meta)
}

func resourceNutanixRecoveryPlanDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.DeleteRecoveryPlan(d.Id())

	if err != nil {
		return fmt.Errorf("error deleting protection_rule id %s): %s", d.Id(), err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for recovery_plan (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixRecoveryPlanExists(conn *v3.Client, name string) (*string, error) {
	var uuid *string

	protectionList, err := conn.V3.ListAllRecoveryPlans()

	if err != nil {
		return nil, err
	}

	for _, protection := range protectionList.Entities {
		if protection.Status.Name == name {
			uuid = protection.Metadata.UUID
		}
	}
	return uuid, nil
}

func getRecoveryPlanResources(d *schema.ResourceData, rp *v3.RecoveryPlanResources) {
	rp.StageList = expandStageList(d)
	rp.Parameters = expandParameters(d)
}

func expandStageList(d *schema.ResourceData) []*v3.StageList {
	stageList := make([]*v3.StageList, 0)
	if v, ok := d.GetOk("stage_list"); ok {
		zoneList := v.([]interface{})
		for _, zone := range zoneList {
			sl := &v3.StageList{}
			v1 := zone.(map[string]interface{})
			if v2, ok := v1["stage_uuid"]; ok && v2.(string) != "" {
				sl.StageUUID = v2.(string)
			}
			if v2, ok := v1["delay_time_secs"]; ok {
				sl.DelayTimeSecs = utils.Int64Ptr(cast.ToInt64(v2))
			}
			if v2, ok := v1["stage_work"].([]interface{}); ok && len(v2) > 0 {
				sl.StageWork = expandStageWork(v2[0].(map[string]interface{}))
			}
			stageList = append(stageList, sl)
		}
	}

	return stageList
}

func expandStageWork(d map[string]interface{}) *v3.StageWork {
	sw := &v3.StageWork{}
	if v1, ok := d["recover_entities"].([]interface{}); ok && len(v1) > 0 {
		recoverEntities := &v3.RecoverEntities{}
		v2 := v1[0].(map[string]interface{})
		if v4, ok := v2["entity_info_list"].([]interface{}); ok && len(v4) > 0 {
			recoverEntities.EntityInfoList = expandEntityInfoList(v4)
		}
		sw.RecoverEntities = recoverEntities
	}
	return sw
}

func expandEntityInfoList(d []interface{}) []*v3.EntityInfoList {
	entities := make([]*v3.EntityInfoList, 0)
	for _, val := range d {
		v := val.(map[string]interface{})
		entity := &v3.EntityInfoList{}
		if v1, ok1 := v["any_entity_reference_kind"]; ok1 && v1.(string) != "" {
			entity.AnyEntityReference.Kind = utils.StringPtr(v1.(string))
		}
		if v1, ok1 := v["any_entity_reference_uuid"]; ok1 && v1.(string) != "" {
			entity.AnyEntityReference.UUID = utils.StringPtr(v1.(string))
		}
		if v1, ok1 := v["any_entity_reference_name"]; ok1 && v1.(string) != "" {
			entity.AnyEntityReference.Name = utils.StringPtr(v1.(string))
		}
		if v1, ok1 := v["categories"]; ok1 {
			entity.Categories = expandCategories(v1)
		}

		entities = append(entities, entity)
	}
	return entities
}

func expandParameters(d *schema.ResourceData) *v3.Parameters {
	parameter := &v3.Parameters{}
	if v, ok := d.GetOk("parameters"); ok {
		v := v.([]interface{})
		for _, v := range v {
			v1 := v.(map[string]interface{})
			if v1, ok1 := v1["floating_ip_assignment_list"].([]interface{}); ok1 {
				parameter.FloatingIPAssignmentList = expandFloatingAssignmentList(v1)
			}
			if v1, ok1 := v1["network_mapping_list"]; ok1 {
				list := v1.([]interface{})
				networkMappingList := make([]*v3.NetworkMappingList, 0)
				networkMapping := &v3.NetworkMappingList{}
				for _, network := range list {
					v2 := network.(map[string]interface{})
					if v2, ok1 := v2["availability_zone_network_mapping_list"].([]interface{}); ok1 {
						networkMapping.AvailabilityZoneNetworkMappingList = expandZoneNetworkMappingList(v2)
					}
					networkMappingList = append(networkMappingList, networkMapping)
				}

				parameter.NetworkMappingList = networkMappingList
			}
		}
	}
	return parameter
}

func expandFloatingAssignmentList(d []interface{}) []*v3.FloatingIPAssignmentList {
	floatings := make([]*v3.FloatingIPAssignmentList, 0)
	for _, float := range d {
		floating := &v3.FloatingIPAssignmentList{}
		v1 := float.(map[string]interface{})
		if v2, ok1 := v1["availability_zone_url"]; ok1 && v2.(string) != "" {
			floating.AvailabilityZoneURL = v2.(string)
		}
		if v2, ok1 := v1["vm_ip_assignment_list"].([]interface{}); ok1 {
			floating.VMIPAssignmentList = expandVMIPAssignmentList(v2)
			floatings = append(floatings, floating)
		}
	}
	return floatings
}

func expandVMIPAssignmentList(d []interface{}) []*v3.VMIPAssignmentList {
	assigns := make([]*v3.VMIPAssignmentList, 0)
	for _, assignment := range d {
		vmial := &v3.VMIPAssignmentList{}
		v1 := assignment.(map[string]interface{})
		if v2, ok1 := v1["test_floating_ip_config"]; ok1 {
			v4 := v2.([]interface{})
			for _, v6 := range v4 {
				v7 := v6.(map[string]interface{})
				if v5, ok1 := v7["ip"]; ok1 && v5.(string) != "" {
					vmial.TestFloatingIPConfig.IP = v5.(string)
				}
				if v5, ok1 := v7["should_allocate_dynamically"]; ok1 {
					vmial.TestFloatingIPConfig.ShouldAllocateDynamically = utils.BoolPtr(v5.(bool))
				}
			}
		}
		if v2, ok1 := v1["recovery_floating_ip_config"]; ok1 {
			v4 := v2.([]interface{})
			for _, v6 := range v4 {
				v7 := v6.(map[string]interface{})
				if v5, ok1 := v7["ip"]; ok1 && v5.(string) != "" {
					vmial.TestFloatingIPConfig.IP = v5.(string)
				}
				if v5, ok1 := v7["should_allocate_dynamically"]; ok1 {
					vmial.TestFloatingIPConfig.ShouldAllocateDynamically = utils.BoolPtr(v5.(bool))
				}
			}
		}
		if v2, ok1 := v1["vm_reference"]; ok1 {
			v6 := v2.([]interface{})
			for _, v7 := range v6 {
				v4 := v7.(map[string]interface{})
				if v5, ok1 := v4["name"]; ok1 && v5.(string) != "" {
					vmial.VMReference.Name = utils.StringPtr(v5.(string))
				}
				if v5, ok1 := v4["uuid"]; ok1 && v5.(string) != "" {
					vmial.VMReference.UUID = utils.StringPtr(v5.(string))
				}
				if v5, ok1 := v4["kind"]; ok1 && v5.(string) != "" {
					vmial.VMReference.Kind = utils.StringPtr(v5.(string))
				}
			}
		}
		if v2, ok1 := v1["vm_nic_information"]; ok1 {
			v6 := v2.([]interface{})
			for _, v7 := range v6 {
				v4 := v7.(map[string]interface{})
				if v5, ok1 := v4["ip"]; ok1 && v5.(string) != "" {
					vmial.VMNICInformation.IP = v5.(string)
				}
				if v5, ok1 := v4["uuid"]; ok1 && v5.(string) != "" {
					vmial.VMNICInformation.UUID = v5.(string)
				}
			}
		}
		assigns = append(assigns, vmial)
	}
	return assigns
}

func expandZoneNetworkMappingList(d []interface{}) []*v3.AvailabilityZoneNetworkMappingList {
	mapping := make([]*v3.AvailabilityZoneNetworkMappingList, 0)
	for _, networkMap := range d {
		netMap := &v3.AvailabilityZoneNetworkMappingList{}
		v4 := networkMap.(map[string]interface{})
		if v5, ok1 := v4["availability_zone_url"]; ok1 && v5.(string) != "" {
			netMap.AvailabilityZoneURL = v5.(string)
		}
		if v5, ok1 := v4["recovery_network"].([]interface{}); ok1 && len(v5) > 0 {
			netMap.RecoveryNetwork = expandRecoveryNetwork(v5[0].([]interface{}))
		}
		if v5, ok1 := v4["test_network"].([]interface{}); ok1 && len(v5) > 0 {
			netMap.RecoveryNetwork = expandRecoveryNetwork(v5[0].([]interface{}))
		}
		mapping = append(mapping, netMap)
	}
	return mapping
}

func expandRecoveryNetwork(d []interface{}) *v3.Network {
	network := &v3.Network{}
	for _, v1 := range d {
		d := v1.(map[string]interface{})
		if v, ok1 := d["name"]; ok1 && v.(string) != "" {
			network.Name = v.(string)
		}
		if v, ok1 := d["virtual_network_reference"].([]interface{}); ok1 {
			network.VirtualNetworkReference = validateRefList(v, utils.StringPtr("recovery_plan"))
		}
		if v, ok1 := d["vpc_reference"].([]interface{}); ok1 {
			network.VPCReference = validateRefList(v, utils.StringPtr("recovery_plan"))
		}
		if v, ok1 := d["name"]; ok1 && v.(string) != "" {
			network.Name = v.(string)
		}
		if v, ok1 := d["subnet_list"].([]interface{}); ok1 {
			network.SubnetList = expandSubnetList(v)
		}
		if v, ok1 := d["use_vpc_reference"].(bool); ok1 {
			network.UseVPCReference = utils.BoolPtr(v)
		}
	}

	return network
}

func expandSubnetList(d []interface{}) []*v3.SubnetList {
	subnets := make([]*v3.SubnetList, 0)
	for _, subnet := range d {
		sub := &v3.SubnetList{}
		v2 := subnet.(map[string]interface{})
		if v3, ok1 := v2["gateway_ip"]; ok1 && v3.(string) != "" {
			sub.GatewayIP = v3.(string)
		}
		if v3, ok1 := v2["external_connectivity_state"]; ok1 && v3.(string) != "" {
			sub.ExternalConnectivityState = v3.(string)
		}
		if v3, ok1 := v2["prefix_length"]; ok1 {
			sub.PrefixLength = utils.Int64Ptr(cast.ToInt64(v3))
		}
		subnets = append(subnets, sub)
	}
	return subnets
}

func flattenStageList(sl []*v3.StageList) []map[string]interface{} {
	stageList := make([]map[string]interface{}, 0)
	for _, v := range sl {
		stage := make(map[string]interface{})

		stage["stage_uuid"] = v.StageUUID
		stage["delay_time_secs"] = utils.Int64Value(v.DelayTimeSecs)
		stage["stage_work"] = flattenStageWork(v.StageWork)

		stageList = append(stageList, stage)
	}
	return stageList
}

func flattenStageWork(stageWork *v3.StageWork) []interface{} {
	sw := make([]interface{}, 0)
	if stageWork.RecoverEntities != nil {
		recoverEntities := make(map[string]interface{})
		recoverEntities["recover_entities"] = flattenEntityInfoList(stageWork.RecoverEntities.EntityInfoList)
		sw = append(sw, recoverEntities)
	}
	return sw
}

func flattenEntityInfoList(entitiesList []*v3.EntityInfoList) []interface{} {
	entities := make([]interface{}, 0)
	log.Printf("[DEBUG] entitiesList: %+v", entitiesList)
	utils.PrintToJSON(entitiesList, "")
	for _, v2 := range entitiesList {
		ent := make(map[string]interface{})
		if v2 != nil {
			entity := make(map[string]interface{})
			entList := make([]interface{}, 0)
			if v2.AnyEntityReference != nil {
				entity["any_entity_reference_name"] = utils.StringValue(v2.AnyEntityReference.Name)
				entity["any_entity_reference_uuid"] = utils.StringValue(v2.AnyEntityReference.UUID)
				entity["any_entity_reference_kind"] = utils.StringValue(v2.AnyEntityReference.Kind)
			}
			entity["categories"] = flattenCategories(v2.Categories)
			entList = append(entList, entity)
			ent["entity_info_list"] = entList
		}
		entities = append(entities, ent)
	}
	log.Printf("[DEBUG] flattenEntityInfoList result: %+v", entities)
	return entities
}

func flattenParameters(par *v3.Parameters) []interface{} {
	parameters := make([]interface{}, 0)
	if par != nil {
		parameter := make(map[string]interface{})
		parameter["floating_ip_assignment_list"] = flattenFloatingAssignmentList(par.FloatingIPAssignmentList)
		parameter["network_mapping_list"] = flattenNetworkMappingList(par.NetworkMappingList)
		parameters = append(parameters, parameter)
	}

	return parameters
}

func flattenFloatingAssignmentList(floatingList []*v3.FloatingIPAssignmentList) []map[string]interface{} {
	floatings := make([]map[string]interface{}, 0)
	if len(floatingList) > 0 {
		for _, floating := range floatingList {
			float := make(map[string]interface{})
			float["availability_zone_url"] = floating.AvailabilityZoneURL
			float["vm_ip_assignment_list"] = flattenVMAssignmentList(floating.VMIPAssignmentList)
			floatings = append(floatings, float)
		}
	}
	return floatings
}

func flattenVMAssignmentList(vmList []*v3.VMIPAssignmentList) []map[string]interface{} {
	assignments := make([]map[string]interface{}, 0)
	if len(vmList) > 0 {
		for _, assignment := range vmList {
			assign := make(map[string]interface{})
			assign["vm_reference"] = flattenReferenceValues(assignment.VMReference)
			floatingConfig := make(map[string]interface{})
			floatingConfig["ip"] = assignment.TestFloatingIPConfig.IP
			floatingConfig["should_allocate_dynamically"] = utils.BoolValue(assignment.TestFloatingIPConfig.ShouldAllocateDynamically)
			assign["test_floating_ip_config"] = floatingConfig
			floatingConfig = make(map[string]interface{})
			floatingConfig["ip"] = assignment.RecoveryFloatingIPConfig.IP
			floatingConfig["should_allocate_dynamically"] = utils.BoolValue(assignment.RecoveryFloatingIPConfig.ShouldAllocateDynamically)
			assign["recovery_floating_ip_config"] = floatingConfig
			assignments = append(assignments, assign)
		}
	}
	return assignments
}

func flattenNetworkMappingList(networksList []*v3.NetworkMappingList) []map[string]interface{} {
	networks := make([]map[string]interface{}, 0)
	if len(networksList) > 0 {
		for _, network := range networksList {
			availibility := make(map[string]interface{})
			zones := make([]map[string]interface{}, 0)
			if len(network.AvailabilityZoneNetworkMappingList) > 0 {
				for _, zone := range network.AvailabilityZoneNetworkMappingList {
					zon := make(map[string]interface{})
					zon["availability_zone_url"] = zone.AvailabilityZoneURL
					zon["recovery_network"] = flattenRecoveryNetwork(zone.RecoveryNetwork)
					zon["test_network"] = flattenRecoveryNetwork(zone.RecoveryNetwork)
					zones = append(zones, zon)
				}
				availibility["availability_zone_network_mapping_list"] = zones
			}
			networks = append(networks, availibility)
		}
	}
	return networks
}

func flattenRecoveryNetwork(d *v3.Network) []interface{} {
	networks := make([]interface{}, 0)
	network := make(map[string]interface{})
	network["name"] = d.Name
	network["virtual_network_reference"] = flattenReferenceValuesList(d.VirtualNetworkReference)
	network["vpc_reference"] = flattenReferenceValuesList(d.VPCReference)
	network["subnet_list"] = flattenSubnetList(d.SubnetList)
	network["use_vpc_reference"] = d.UseVPCReference
	networks = append(networks, network)
	return networks
}

func flattenSubnetList(subnets []*v3.SubnetList) []map[string]interface{} {
	subs := make([]map[string]interface{}, 0)
	if len(subnets) > 0 {
		for _, subnet := range subnets {
			sub := make(map[string]interface{})
			sub["gateway_ip"] = subnet.GatewayIP
			sub["external_connectivity_state"] = subnet.ExternalConnectivityState
			sub["prefix_length"] = utils.Int64Value(subnet.PrefixLength)
			subs = append(subs, sub)
		}
	}
	return subs
}
