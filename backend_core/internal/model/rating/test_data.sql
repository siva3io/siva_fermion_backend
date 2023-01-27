SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET session_replication_role = 'replica';
SET schema 'public';
/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/

INSERT INTO public.rating_categories (id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, category_name, rating_value) VALUES (1, true, true, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'ITEM', 5);
INSERT INTO public.rating_categories (id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, category_name, rating_value) VALUES (2, true, true, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'PROVIDER', 5);
SELECT setval( pg_get_serial_sequence('public.rating_categories', 'id'), (SELECT max(id) FROM public.rating_categories));

INSERT INTO public.feedback_categories(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,rating_category_name)VALUES(1, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,'ITEM');
INSERT INTO public.feedback_categories(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,rating_category_name)VALUES(2, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,'PROVIDER');
SELECT setval( pg_get_serial_sequence('public.feedback_categories', 'id'), (SELECT max(id) FROM public.feedback_categories));

INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(3, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,3,3,'{}',1);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(4, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,4,1,'{}',1);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(1, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,5,'{}',1);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(5, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,5,2,'{}',1);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(6, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,6,5,'{}',2);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(2, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,2,4,'{}',2);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(7, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,7,4,'{}',2);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(8, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,8,3,'{}',2);
INSERT INTO public.ratings(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,related_id,rating_value,feed_back,rating_category_id)VALUES(9, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,9,2,'{}',2);
SELECT setval( pg_get_serial_sequence('public.ratings', 'id'), (SELECT max(id) FROM public.ratings));

INSERT INTO public.feedback_forms(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_category_id,rating_value,linked_url,tl_method,additional_prop1,additional_prop2,additional_prop3)VALUES(1, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30',1, NULL,1,5,'eunimart.com','post','property1','property2','property3');
INSERT INTO public.feedback_forms(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_category_id,rating_value,linked_url,tl_method,additional_prop1,additional_prop2,additional_prop3)VALUES(2, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,2,5,'eunimart.com','put','property1','property2','property3');
INSERT INTO public.feedback_forms(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_category_id,rating_value,linked_url,tl_method,additional_prop1,additional_prop2,additional_prop3)VALUES(3, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,4,'eunimart.com','post','property1','property2','property3');
INSERT INTO public.feedback_forms(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_category_id,rating_value,linked_url,tl_method,additional_prop1,additional_prop2,additional_prop3)VALUES(4, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,2,4,'eunimart.com','post','property1','property2','property3');
SELECT setval( pg_get_serial_sequence('public.feedback_forms', 'id'), (SELECT max(id) FROM public.feedback_forms));

INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(1, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,1,'How is the product',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(2, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 2, NULL,1,1,'How is the shipping agents behaviour?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(3, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 2, NULL,1,1,'How is your experience with us?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(4, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,1,'How likely are you to recommend this product to others?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(5, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 2, NULL,1,1,'How could we improve our product to better meet your needs?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(6, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,1,'What, if any, products, services, or features are we missing?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(7, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 2, NULL,1,1,'Where did you first hear about us?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
INSERT INTO public.feedbacks(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,feedback_form_id,parent_id,question,answer_type_id)VALUES(8, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,1,1,'Have you used a similar [product or service] before?',(SELECT public.lookupcodes.id FROM public.lookuptypes,public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'ANSWER_TYPE' AND public.lookupcodes.lookup_code = 'TEXT'));
SELECT setval( pg_get_serial_sequence('public.feedbacks', 'id'), (SELECT max(id) FROM public.feedbacks));


INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(2, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(3, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(4, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(5, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(6, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(7, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(8, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(9, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(10, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,' {
    "context": {
        "domain": "nic2004:52110",
        "country": "IND",
        "city": "std:080",
        "action": "issue",
        "core_version": "1.0.0",
        "bap_id": "ondc.bap.com",
        "bap_uri": "https://abc.interfacing.app.com",
        "bpp_id": "ondc.bpp.com",
        "bpp_uri": "https://abc.transaction.counterpary.app.com",
        "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
        "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
        "timestamp": "2022-10-28T10:34:58.469Z",
        "ttl": "PT30S"
    },
    "message": {
        "issue": {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "complainant_info": {
                "person": {
                    "name": "Sam Manuel"
                },
                "contact": {
                    "name": "Sam Manuel",
                    "address": {
                        "full": "45, Park Avenue, Delhi",
                        "format": "string"
                    },
                    "phone": "9879879870",
                    "email": "sam@yahoo.com"
                },
                "complainant_action": "Open",
                "complainant_action_reason": "identified an issue with order so opening an issue"
            },
            "order": {
                "id": "string",
                "ref_order_ids": [
                    "1209",
                    "1210"
                ],
                "items": [
                    {
                        "id": "23433",
                        "parent_item_id": "44334",
                        "descriptor": {
                            "name": "laptop",
                            "code": "code",
                            "short_desc": "laptop description",
                            "long_desc": "laptop long description",
                            "additional_desc": {
                                "url": "string",
                                "content_type": "text/plain"
                            }
                        },
                        "manufacturer": {
                            "descriptor": {
                                "name": "manuf desc",
                                "code": "mfg93",
                                "short_desc": "mfg desc ",
                                "long_desc": "mfg long desc"
                            },
                            "address": {
                                "full": "433054, Park Avenue, Delhi",
                                "format": "string"
                            },
                            "contact": {
                                "name": "Sam X",
                                "address": {
                                    "full": "4993, Park Avenue ,Delhi",
                                    "format": "string"
                                },
                                "phone": "4059404549",
                                "email": "email@email.com"
                            }
                        },
                        "price": {
                            "currency": "INR",
                            "value": "590"
                        },
                        "quantity": {
                            "allocated": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "3",
                                    "unit": "string"
                                }
                            },
                            "available": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "6",
                                    "unit": "string"
                                }
                            },
                            "maximum": {
                                "count": 1,
                                "measure": {
                                    "type": "CONSTANT",
                                    "value": "2",
                                    "unit": "string"
                                }
                            }
                        },
                        "category_id": "409",
                        "fulfillment_id": "230"
                    }
                ],
                "fulfillments": [
                    {
                        "id": "230",
                        "type": "xyz",
                        "provider_id": "30493",
                        "state": {
                            "descriptor": {
                                "name": "fulfillments state name",
                                "code": "fulfillments state code",
                                "short_desc": "fulfillments desc",
                                "long_desc": "fulfillments long desc name"
                            },
                            "updated_at": "2022-12-08T05:5:34.946Z",
                            "updated_by": "string"
                        },
                        "tracking": false,
                        "customer": {
                            "person": {
                                "id": "101",
                                "name": "Sam B"
                            },
                            "contact": {
                                "name": "Sam B",
                                "address": {
                                    "full": "43450, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "5904500590",
                                "email": "sam@emial.com"
                            }
                        },
                        "agent": {
                            "person": {
                                "id": "40904",
                                "name": "Sam A"
                            },
                            "contact": {
                                "name": "Sam F C",
                                "address": {
                                    "full": "54059, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "phone": "4304903440",
                                "email": "email@pa.com"
                            },
                            "organization": {
                                "descriptor": {
                                    "name": "FAOrg",
                                    "code": "09403",
                                    "short_desc": "FA ORG shot desc",
                                    "long_desc": "FA Logn shot desc"
                                },
                                "address": {
                                    "full": "494589, Park Avenue, Delhi",
                                    "format": "string"
                                },
                                "contact": {
                                    "name": "Sam FAO",
                                    "address": {
                                        "full": "95904, Park Avenue Delhi",
                                        "format": "string"
                                    },
                                    "phone": "9403904940",
                                    "email": "fao@email.com"
                                }
                            }
                        },
                        "contact": {
                            "name": "con",
                            "address": {
                                "full": "4390, Park Avenue Delhi",
                                "format": "string"
                            },
                            "phone": "0549505400",
                            "email": "dmai@email.com"
                        }
                    }
                ],
                "created_at": "2022-12-08T05:5:34.961Z",
                "updated_at": "2022-12-08T0:51:34.961Z"
            },
            "description": {
                "name": "issue is with the product",
                "short_desc": "product is duplicate",
                "long_desc": "the product appears to be of bad quality",
                "additional_desc": {
                    "url": "https://link-to-assests.buyer.com/product-imag.img",
                    "content_type": "text/plain"
                }
            },
            "category": "Product",
            "issue_type": "Issue",
            "source": {
                "id": "813389c0-65a2-11ed-9022-0242ac120002",
                "issue_source_type": "Consumer"
            },
            "supplementary_information": [
                {
                    "issue_update_info": {
                        "short_desc": "",
                        "long_desc": ""
                    },
                    "created_at": "",
                    "modified_at": ""
                }
            ],
            "expected_response_time": {
                "duration": "T1H"
            },
            "expected_resolution_time": {
                "duration": "P1D"
            },
            "status": {
                "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "status": "Open",
                "closing_reason": "closing issue since the issue is resolved",
                "closing_remarks": "closed issue on resolution",
                "status_details": {
                    "short_desc": "issue status details",
                    "long_desc": "issue status long description"
                },
                "status_change_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "issue_modification_date": {
                    "timestamp": "2022-12-04T05:28:11.320Z"
                },
                "modified_by": {
                    "org": {
                        "name": "org name of status modifier"
                    },
                    "contact": {
                        "name": "contact name of status modifier",
                        "address": {
                            "full": "40595, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "04950394039",
                        "email": "modiemail@org.com"
                    },
                    "person": {
                        "name": "person name of status modifier"
                    }
                },
                "modified_issue_field_name": "string"
            },
            "created_at": "2022-121.320Z",
            "modified_at": "2022-111.320Z",
            "finalized_odr": {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization": {
                    "org": {
                        "name": "ODR provider"
                    },
                    "contact": {
                        "name": "ODR provider contact name",
                        "address": {
                            "full": "A-3040, park Avenue Delhi",
                            "format": "string"
                        },
                        "phone": "0495049304",
                        "email": "email@odr.com"
                    },
                    "person": {
                        "name": "odr provider contact  person name"
                    }
                },
                "pricing_model": {
                    "price": {
                        "currency": "string",
                        "value": "10000"
                    },
                    "pricing_info": "10% of issue cost"
                },
                "resolution_ratings": {
                    "rating_value": "98%"
                }
            },
            "rating": "THUMBS-UP"
        }
    }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(11, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,
'{
  "context":
  {
    "domain": "nic2004:52110",
    "country": "IND",
    "city": "std:080",
    "action": "on_issue",
    "core_version": "1.0.0",
    "bap_id": "ondc.bap.com",
    "bap_uri": "https://abc.interfacing.app.com",
    "bpp_id": "ondc.bpp.com",
    "bpp_uri": "https://abc.transaction.counterpary.app.com",
    "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
    "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
    "timestamp": "2022-10-28T10:34:58.469Z"
  },
  "message":
  {
    "issue_resolution_details":
    {
      "issue":
      {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "complainant_info":
        {
          "person":
          {
            "name": "Sam Manuel"
          },
          "contact":
          {
            "name": "Sam Manuel",
            "address":
            {
              "full": "45,Park Avenue, Delhi",
              "format": "string"
            },
            "phone": "9879879870",
            "email": "sam@yahoo.com"
          },
          "complainant_action": "open",
          "complainant_action_reason": "identified an issue with order so opening an issue"
        },
        "order":
        {
          "id": "0949848944",
          "ref_order_ids":
          [
            "1209",
            "1210"
          ],
          "items":
          [ {
              "id": "23433",
              "parent_item_id": "44334",
              "descriptor":
              {
                "name": "laptop",
                "code": "code",
                "short_desc": "laptop description",
                "long_desc": "laptop long description",
                "additional_desc":
                {
                  "url": "string",
                  "content_type": "text/plain"
                }
              },
              "manufacturer":
              {
                "descriptor":
                {
                  "name": "manuf desc",
                  "code": "mfg93",
                  "short_desc": "mfg desc ",
                  "long_desc": "mfg long desc"
                },
                "address":
                {
                  "full": "433054, Park Avenue, Delhi",
                  "format": "string"
                },
                "contact":
                {
                  "name": "Sam X",
                  "address":
                  {
                    "full": "4993, Park Avenue ,Delhi",
                    "format": "string"
                  },
                  "phone": "4059404549",
                  "email": "email@email.com"
                }
              },
              "price":
              {
                "currency": "INR",
                "value": "590"
              },
              "quantity":
              {
                "allocated":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "3",
                    "unit": "string"
                  }
                },
                "available":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "6",
                    "unit": "string"
                  }
                },
                "maximum":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "2",
                    "unit": "string"
                  }
                }
              },
              "category_id": "409",
              "fulfillment_id": "230"
            }
          ],
          "fulfillments":
          [ {
              "id": "230",
              "type": "xyz",
              "provider_id": "30493",
              "state":
              {
                "descriptor":
                {
                  "name": "fulfillments state name",
                  "code": "fulfillments state code",
                  "short_desc": "fulfillments desc",
                  "long_desc": "fulfillments long desc name"
                },
                "updated_at": "2022-12-08T05:51:34.946Z",
                "updated_by": "string"
              },
              "tracking": false,
              "customer":
              {
                "person":
                {
                  "id": "101",
                  "name": "Sam B"
                },
                "contact":
                {
                  "name": "Sam B",
                  "address":
                  {
                    "full": "43450, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "5904500590",
                  "email": "sam@emial.com"
                }
              },
              "agent":
              {
                "person":
                {
                  "id": "40904",
                  "name": "Sam A"
                },
                "contact":
                {
                  "name": "Sam F C",
                  "address":
                  {
                    "full": "54059, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "4304903440",
                  "email": "email@pa.com"
                },
                "organization":
                {
                  "descriptor":
                  {
                    "name": "FAOrg",
                    "code": "09403",
                    "short_desc": "FA ORG shot desc",
                    "long_desc": "FA Logn shot desc"
                  },
                  "address":
                  {
                    "full": "494589, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "contact":
                  {
                    "name": "Sam FAO",
                    "address":
                    {
                      "full": "95904, Park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "9403904940",
                    "email": "fao@email.com"
                  }
                }
              },
              "contact":
              {
                "name": "con",
                "address":
                {
                  "full": "4390, Park Avenue Delhi",
                  "format": "string"
                },
                "phone": "0549505400",
                "email": "dmai@email.com"
              }
            }
          ],
          "created_at": "2022-12-08T05:51:34.961Z",
          "updated_at": "2022-12-08T05:51:34.961Z"
        },
        "description":
        {
          "name": "issue is with the product",
          "short_desc": "product is duplicate",
          "long_desc": "the product appears to be of bad quality",
          "additional_desc":
          {
            "url": "https://link-to-assests.buyer.com/product-imag.img",
            "content_type": "text/plain"
          }
        },
        "category": "Product",
        "issue_type": "Issue",
        "source":
        {
          "network_participant_id": "813389c0-65a2-11ed-9022-0242ac120002",
          "issue_source_type": "Consumer"
        },
        "supplementary_information":
        [ {
            "issue_update_info":
            {
              "name": "string",
              "code": "respondent",
              "short_desc": "provide further proofs of the issue",
              "long_desc": "provide some pictures to verify if the product is actually not of the right quality"
            },
            "created_at": "2022-12-04T06:30:47.753Z",
            "modified_at": "2022-12-04T06:30:47.753Z"
          }
        ],
        "expected_response_time":
        {
          "duration": "T1H"
        },
        "expected_resolution_time":
        {
          "duration": "P1D"
        },
        "status":
        {
          "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "status": "Resolved",
          "status_details":
          {
            "code": "Accepted",
            "short_desc": "issue status details ",
            "long_desc": "issue status long description"
          },
          "status_change_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "issue_modification_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "modified_by":
          {
            "org":
            {
              "name": "org name of status modifier"
            },
            "contact":
            {
              "name": "contact name of status modifier",
              "address":
              {
                "full": "40595, park Avenue Delhi",
                "format": "string"
              },
              "phone": "04950394039",
              "email": "modiemail@org.com"
            },
            "person":
            {
              "name": "person name of status modifier"
            }
          },
          "modified_issue_field_name": "string"
        },
        "created_at": "2022-12-04T06:30:47.753Z",
        "modified_at": "2022-12-04T06:30:47.753Z"
      },
      "resolution_provider":
      {
        "respondent_info":
        {
          "type": "Transaction Counterparty NP",
          "organization":
          {
            "org":
            {
              "name": "interfacing app resolution provider app name"
            },
            "contact":
            {
              "name": "Sam C",
              "address":
              {
                "full": "4959, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "4059304940",
              "email": "email@resolutionproviderorg.com"
            },
            "person":
            {
              "name": "resolution provider org contact person name"
            }
          },
          "resolution_support":
          {
            "respondentChatLink": "http://chat-link/respondant",
            "respondentEmail": "email@respondant.com",
            "respondentContact":
            {
              "name": "resolution providers, issue resolution support respondent name",
              "address":
              {
                "full": "4955, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "5949595059",
              "email": "respondantemail@resolutionprovider.com"
            },
            "respondentFaqs":
            {
              "additionalProp1":
              {
                "question": "Question 1",
                "Answer": "Answer 1"
              },
              "additionalProp2":
              {
                "question": "question 2",
                "Answer": "answer 2"
              },
              "additionalProp3":
              {
                "question": "Question 3",
                "Answer": "Answer 3"
              }
            },
            "additional_sources":
            {
              "additionalProp1":
              {
                "type": "additional resolution source type : eg : IVR",
                "link": "ivr contact details or link to it : eg: http://resolutionsupport.ivr.resolutionprovider.com"
              }
            },
            "gros":
            [ {
              "person": {
                "name": "Sam D",
                "image": "https://gro/profile-image.png"
              },
              "contact":
              {
                "name": "Sam D",
                "address":
                {
                  "full": "45 park anv. Delhi",
                  "format": "string"
                },
                "phone": "9605960796",
                "email": "email@gro.com"
              },
              "gro_type": "Interfacing NP GRO"
              }
          ],
            "selected_odrs":
            [ {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization":
                {
                  "org":
                  {
                    "name": "ODR provider"
                  },
                  "contact":
                  {
                    "name": "ODR provider contact name",
                    "address":
                    {
                      "full": "A-3040, park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "0495049304",
                    "email": "email@odr.com"
                  },
                  "person":
                  {
                    "name": "odr provider contact  person name"
                  }
                },
                "pricing_model":
                {
                  "price":
                   {
                    "currency": "rupees",
                    "value": "10000"
                   },
                  "pricing_info": "10% of invoice cost"
                },
                "resolution_ratings":
                {
                  "rating_value": "98%"
                }
              }
            ]
          }
        }
      },
      "resolution": {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "resolution": "issue resolution details",
        "resolution_remarks": "remarks related to the resolution",
        "gro_remarks": "remarks from gro related to issue resolution",
        "dispute_resolution_remarks": "provided by odr about dispute resolution",
        "resolution_action": "resolve"
      }
    }
  },
  "error":
  {
    "code": "string",
    "path": "string",
    "message": "string"
  }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(12, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,
'{
  "context":
  {
    "domain": "nic2004:52110",
    "country": "IND",
    "city": "std:080",
    "action": "on_issue",
    "core_version": "1.0.0",
    "bap_id": "ondc.bap.com",
    "bap_uri": "https://abc.interfacing.app.com",
    "bpp_id": "ondc.bpp.com",
    "bpp_uri": "https://abc.transaction.counterpary.app.com",
    "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
    "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
    "timestamp": "2022-10-28T10:34:58.469Z"
  },
  "message":
  {
    "issue_resolution_details":
    {
      "issue":
      {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "complainant_info":
        {
          "person":
          {
            "name": "Sam Manuel"
          },
          "contact":
          {
            "name": "Sam Manuel",
            "address":
            {
              "full": "45,Park Avenue, Delhi",
              "format": "string"
            },
            "phone": "9879879870",
            "email": "sam@yahoo.com"
          },
          "complainant_action": "open",
          "complainant_action_reason": "identified an issue with order so opening an issue"
        },
        "order":
        {
          "id": "0949848944",
          "ref_order_ids":
          [
            "1209",
            "1210"
          ],
          "items":
          [ {
              "id": "23433",
              "parent_item_id": "44334",
              "descriptor":
              {
                "name": "laptop",
                "code": "code",
                "short_desc": "laptop description",
                "long_desc": "laptop long description",
                "additional_desc":
                {
                  "url": "string",
                  "content_type": "text/plain"
                }
              },
              "manufacturer":
              {
                "descriptor":
                {
                  "name": "manuf desc",
                  "code": "mfg93",
                  "short_desc": "mfg desc ",
                  "long_desc": "mfg long desc"
                },
                "address":
                {
                  "full": "433054, Park Avenue, Delhi",
                  "format": "string"
                },
                "contact":
                {
                  "name": "Sam X",
                  "address":
                  {
                    "full": "4993, Park Avenue ,Delhi",
                    "format": "string"
                  },
                  "phone": "4059404549",
                  "email": "email@email.com"
                }
              },
              "price":
              {
                "currency": "INR",
                "value": "590"
              },
              "quantity":
              {
                "allocated":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "3",
                    "unit": "string"
                  }
                },
                "available":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "6",
                    "unit": "string"
                  }
                },
                "maximum":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "2",
                    "unit": "string"
                  }
                }
              },
              "category_id": "409",
              "fulfillment_id": "230"
            }
          ],
          "fulfillments":
          [ {
              "id": "230",
              "type": "xyz",
              "provider_id": "30493",
              "state":
              {
                "descriptor":
                {
                  "name": "fulfillments state name",
                  "code": "fulfillments state code",
                  "short_desc": "fulfillments desc",
                  "long_desc": "fulfillments long desc name"
                },
                "updated_at": "2022-12-08T05:51:34.946Z",
                "updated_by": "string"
              },
              "tracking": false,
              "customer":
              {
                "person":
                {
                  "id": "101",
                  "name": "Sam B"
                },
                "contact":
                {
                  "name": "Sam B",
                  "address":
                  {
                    "full": "43450, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "5904500590",
                  "email": "sam@emial.com"
                }
              },
              "agent":
              {
                "person":
                {
                  "id": "40904",
                  "name": "Sam A"
                },
                "contact":
                {
                  "name": "Sam F C",
                  "address":
                  {
                    "full": "54059, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "4304903440",
                  "email": "email@pa.com"
                },
                "organization":
                {
                  "descriptor":
                  {
                    "name": "FAOrg",
                    "code": "09403",
                    "short_desc": "FA ORG shot desc",
                    "long_desc": "FA Logn shot desc"
                  },
                  "address":
                  {
                    "full": "494589, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "contact":
                  {
                    "name": "Sam FAO",
                    "address":
                    {
                      "full": "95904, Park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "9403904940",
                    "email": "fao@email.com"
                  }
                }
              },
              "contact":
              {
                "name": "con",
                "address":
                {
                  "full": "4390, Park Avenue Delhi",
                  "format": "string"
                },
                "phone": "0549505400",
                "email": "dmai@email.com"
              }
            }
          ],
          "created_at": "2022-12-08T05:51:34.961Z",
          "updated_at": "2022-12-08T05:51:34.961Z"
        },
        "description":
        {
          "name": "issue is with the product",
          "short_desc": "product is duplicate",
          "long_desc": "the product appears to be of bad quality",
          "additional_desc":
          {
            "url": "https://link-to-assests.buyer.com/product-imag.img",
            "content_type": "text/plain"
          }
        },
        "category": "Product",
        "issue_type": "Issue",
        "source":
        {
          "network_participant_id": "813389c0-65a2-11ed-9022-0242ac120002",
          "issue_source_type": "Consumer"
        },
        "supplementary_information":
        [ {
            "issue_update_info":
            {
              "name": "string",
              "code": "respondent",
              "short_desc": "provide further proofs of the issue",
              "long_desc": "provide some pictures to verify if the product is actually not of the right quality"
            },
            "created_at": "2022-12-04T06:30:47.753Z",
            "modified_at": "2022-12-04T06:30:47.753Z"
          }
        ],
        "expected_response_time":
        {
          "duration": "T1H"
        },
        "expected_resolution_time":
        {
          "duration": "P1D"
        },
        "status":
        {
          "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "status": "Resolved",
          "status_details":
          {
            "code": "Accepted",
            "short_desc": "issue status details ",
            "long_desc": "issue status long description"
          },
          "status_change_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "issue_modification_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "modified_by":
          {
            "org":
            {
              "name": "org name of status modifier"
            },
            "contact":
            {
              "name": "contact name of status modifier",
              "address":
              {
                "full": "40595, park Avenue Delhi",
                "format": "string"
              },
              "phone": "04950394039",
              "email": "modiemail@org.com"
            },
            "person":
            {
              "name": "person name of status modifier"
            }
          },
          "modified_issue_field_name": "string"
        },
        "created_at": "2022-12-04T06:30:47.753Z",
        "modified_at": "2022-12-04T06:30:47.753Z"
      },
      "resolution_provider":
      {
        "respondent_info":
        {
          "type": "Transaction Counterparty NP",
          "organization":
          {
            "org":
            {
              "name": "interfacing app resolution provider app name"
            },
            "contact":
            {
              "name": "Sam C",
              "address":
              {
                "full": "4959, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "4059304940",
              "email": "email@resolutionproviderorg.com"
            },
            "person":
            {
              "name": "resolution provider org contact person name"
            }
          },
          "resolution_support":
          {
            "respondentChatLink": "http://chat-link/respondant",
            "respondentEmail": "email@respondant.com",
            "respondentContact":
            {
              "name": "resolution providers, issue resolution support respondent name",
              "address":
              {
                "full": "4955, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "5949595059",
              "email": "respondantemail@resolutionprovider.com"
            },
            "respondentFaqs":
            {
              "additionalProp1":
              {
                "question": "Question 1",
                "Answer": "Answer 1"
              },
              "additionalProp2":
              {
                "question": "question 2",
                "Answer": "answer 2"
              },
              "additionalProp3":
              {
                "question": "Question 3",
                "Answer": "Answer 3"
              }
            },
            "additional_sources":
            {
              "additionalProp1":
              {
                "type": "additional resolution source type : eg : IVR",
                "link": "ivr contact details or link to it : eg: http://resolutionsupport.ivr.resolutionprovider.com"
              }
            },
            "gros":
            [ {
              "person": {
                "name": "Sam D",
                "image": "https://gro/profile-image.png"
              },
              "contact":
              {
                "name": "Sam D",
                "address":
                {
                  "full": "45 park anv. Delhi",
                  "format": "string"
                },
                "phone": "9605960796",
                "email": "email@gro.com"
              },
              "gro_type": "Interfacing NP GRO"
              }
          ],
            "selected_odrs":
            [ {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization":
                {
                  "org":
                  {
                    "name": "ODR provider"
                  },
                  "contact":
                  {
                    "name": "ODR provider contact name",
                    "address":
                    {
                      "full": "A-3040, park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "0495049304",
                    "email": "email@odr.com"
                  },
                  "person":
                  {
                    "name": "odr provider contact  person name"
                  }
                },
                "pricing_model":
                {
                  "price":
                   {
                    "currency": "rupees",
                    "value": "10000"
                   },
                  "pricing_info": "10% of invoice cost"
                },
                "resolution_ratings":
                {
                  "rating_value": "98%"
                }
              }
            ]
          }
        }
      },
      "resolution": {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "resolution": "issue resolution details",
        "resolution_remarks": "remarks related to the resolution",
        "gro_remarks": "remarks from gro related to issue resolution",
        "dispute_resolution_remarks": "provided by odr about dispute resolution",
        "resolution_action": "resolve"
      }
    }
  },
  "error":
  {
    "code": "string",
    "path": "string",
    "message": "string"
  }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(13, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,
'{
  "context":
  {
    "domain": "nic2004:52110",
    "country": "IND",
    "city": "std:080",
    "action": "on_issue",
    "core_version": "1.0.0",
    "bap_id": "ondc.bap.com",
    "bap_uri": "https://abc.interfacing.app.com",
    "bpp_id": "ondc.bpp.com",
    "bpp_uri": "https://abc.transaction.counterpary.app.com",
    "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
    "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
    "timestamp": "2022-10-28T10:34:58.469Z"
  },
  "message":
  {
    "issue_resolution_details":
    {
      "issue":
      {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "complainant_info":
        {
          "person":
          {
            "name": "Sam Manuel"
          },
          "contact":
          {
            "name": "Sam Manuel",
            "address":
            {
              "full": "45,Park Avenue, Delhi",
              "format": "string"
            },
            "phone": "9879879870",
            "email": "sam@yahoo.com"
          },
          "complainant_action": "open",
          "complainant_action_reason": "identified an issue with order so opening an issue"
        },
        "order":
        {
          "id": "0949848944",
          "ref_order_ids":
          [
            "1209",
            "1210"
          ],
          "items":
          [ {
              "id": "23433",
              "parent_item_id": "44334",
              "descriptor":
              {
                "name": "laptop",
                "code": "code",
                "short_desc": "laptop description",
                "long_desc": "laptop long description",
                "additional_desc":
                {
                  "url": "string",
                  "content_type": "text/plain"
                }
              },
              "manufacturer":
              {
                "descriptor":
                {
                  "name": "manuf desc",
                  "code": "mfg93",
                  "short_desc": "mfg desc ",
                  "long_desc": "mfg long desc"
                },
                "address":
                {
                  "full": "433054, Park Avenue, Delhi",
                  "format": "string"
                },
                "contact":
                {
                  "name": "Sam X",
                  "address":
                  {
                    "full": "4993, Park Avenue ,Delhi",
                    "format": "string"
                  },
                  "phone": "4059404549",
                  "email": "email@email.com"
                }
              },
              "price":
              {
                "currency": "INR",
                "value": "590"
              },
              "quantity":
              {
                "allocated":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "3",
                    "unit": "string"
                  }
                },
                "available":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "6",
                    "unit": "string"
                  }
                },
                "maximum":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "2",
                    "unit": "string"
                  }
                }
              },
              "category_id": "409",
              "fulfillment_id": "230"
            }
          ],
          "fulfillments":
          [ {
              "id": "230",
              "type": "xyz",
              "provider_id": "30493",
              "state":
              {
                "descriptor":
                {
                  "name": "fulfillments state name",
                  "code": "fulfillments state code",
                  "short_desc": "fulfillments desc",
                  "long_desc": "fulfillments long desc name"
                },
                "updated_at": "2022-12-08T05:51:34.946Z",
                "updated_by": "string"
              },
              "tracking": false,
              "customer":
              {
                "person":
                {
                  "id": "101",
                  "name": "Sam B"
                },
                "contact":
                {
                  "name": "Sam B",
                  "address":
                  {
                    "full": "43450, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "5904500590",
                  "email": "sam@emial.com"
                }
              },
              "agent":
              {
                "person":
                {
                  "id": "40904",
                  "name": "Sam A"
                },
                "contact":
                {
                  "name": "Sam F C",
                  "address":
                  {
                    "full": "54059, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "4304903440",
                  "email": "email@pa.com"
                },
                "organization":
                {
                  "descriptor":
                  {
                    "name": "FAOrg",
                    "code": "09403",
                    "short_desc": "FA ORG shot desc",
                    "long_desc": "FA Logn shot desc"
                  },
                  "address":
                  {
                    "full": "494589, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "contact":
                  {
                    "name": "Sam FAO",
                    "address":
                    {
                      "full": "95904, Park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "9403904940",
                    "email": "fao@email.com"
                  }
                }
              },
              "contact":
              {
                "name": "con",
                "address":
                {
                  "full": "4390, Park Avenue Delhi",
                  "format": "string"
                },
                "phone": "0549505400",
                "email": "dmai@email.com"
              }
            }
          ],
          "created_at": "2022-12-08T05:51:34.961Z",
          "updated_at": "2022-12-08T05:51:34.961Z"
        },
        "description":
        {
          "name": "issue is with the product",
          "short_desc": "product is duplicate",
          "long_desc": "the product appears to be of bad quality",
          "additional_desc":
          {
            "url": "https://link-to-assests.buyer.com/product-imag.img",
            "content_type": "text/plain"
          }
        },
        "category": "Product",
        "issue_type": "Issue",
        "source":
        {
          "network_participant_id": "813389c0-65a2-11ed-9022-0242ac120002",
          "issue_source_type": "Consumer"
        },
        "supplementary_information":
        [ {
            "issue_update_info":
            {
              "name": "string",
              "code": "respondent",
              "short_desc": "provide further proofs of the issue",
              "long_desc": "provide some pictures to verify if the product is actually not of the right quality"
            },
            "created_at": "2022-12-04T06:30:47.753Z",
            "modified_at": "2022-12-04T06:30:47.753Z"
          }
        ],
        "expected_response_time":
        {
          "duration": "T1H"
        },
        "expected_resolution_time":
        {
          "duration": "P1D"
        },
        "status":
        {
          "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "status": "Resolved",
          "status_details":
          {
            "code": "Accepted",
            "short_desc": "issue status details ",
            "long_desc": "issue status long description"
          },
          "status_change_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "issue_modification_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "modified_by":
          {
            "org":
            {
              "name": "org name of status modifier"
            },
            "contact":
            {
              "name": "contact name of status modifier",
              "address":
              {
                "full": "40595, park Avenue Delhi",
                "format": "string"
              },
              "phone": "04950394039",
              "email": "modiemail@org.com"
            },
            "person":
            {
              "name": "person name of status modifier"
            }
          },
          "modified_issue_field_name": "string"
        },
        "created_at": "2022-12-04T06:30:47.753Z",
        "modified_at": "2022-12-04T06:30:47.753Z"
      },
      "resolution_provider":
      {
        "respondent_info":
        {
          "type": "Transaction Counterparty NP",
          "organization":
          {
            "org":
            {
              "name": "interfacing app resolution provider app name"
            },
            "contact":
            {
              "name": "Sam C",
              "address":
              {
                "full": "4959, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "4059304940",
              "email": "email@resolutionproviderorg.com"
            },
            "person":
            {
              "name": "resolution provider org contact person name"
            }
          },
          "resolution_support":
          {
            "respondentChatLink": "http://chat-link/respondant",
            "respondentEmail": "email@respondant.com",
            "respondentContact":
            {
              "name": "resolution providers, issue resolution support respondent name",
              "address":
              {
                "full": "4955, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "5949595059",
              "email": "respondantemail@resolutionprovider.com"
            },
            "respondentFaqs":
            {
              "additionalProp1":
              {
                "question": "Question 1",
                "Answer": "Answer 1"
              },
              "additionalProp2":
              {
                "question": "question 2",
                "Answer": "answer 2"
              },
              "additionalProp3":
              {
                "question": "Question 3",
                "Answer": "Answer 3"
              }
            },
            "additional_sources":
            {
              "additionalProp1":
              {
                "type": "additional resolution source type : eg : IVR",
                "link": "ivr contact details or link to it : eg: http://resolutionsupport.ivr.resolutionprovider.com"
              }
            },
            "gros":
            [ {
              "person": {
                "name": "Sam D",
                "image": "https://gro/profile-image.png"
              },
              "contact":
              {
                "name": "Sam D",
                "address":
                {
                  "full": "45 park anv. Delhi",
                  "format": "string"
                },
                "phone": "9605960796",
                "email": "email@gro.com"
              },
              "gro_type": "Interfacing NP GRO"
              }
          ],
            "selected_odrs":
            [ {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization":
                {
                  "org":
                  {
                    "name": "ODR provider"
                  },
                  "contact":
                  {
                    "name": "ODR provider contact name",
                    "address":
                    {
                      "full": "A-3040, park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "0495049304",
                    "email": "email@odr.com"
                  },
                  "person":
                  {
                    "name": "odr provider contact  person name"
                  }
                },
                "pricing_model":
                {
                  "price":
                   {
                    "currency": "rupees",
                    "value": "10000"
                   },
                  "pricing_info": "10% of invoice cost"
                },
                "resolution_ratings":
                {
                  "rating_value": "98%"
                }
              }
            ]
          }
        }
      },
      "resolution": {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "resolution": "issue resolution details",
        "resolution_remarks": "remarks related to the resolution",
        "gro_remarks": "remarks from gro related to issue resolution",
        "dispute_resolution_remarks": "provided by odr about dispute resolution",
        "resolution_action": "resolve"
      }
    }
  },
  "error":
  {
    "code": "string",
    "path": "string",
    "message": "string"
  }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(14, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,
'{
  "context":
  {
    "domain": "nic2004:52110",
    "country": "IND",
    "city": "std:080",
    "action": "on_issue",
    "core_version": "1.0.0",
    "bap_id": "ondc.bap.com",
    "bap_uri": "https://abc.interfacing.app.com",
    "bpp_id": "ondc.bpp.com",
    "bpp_uri": "https://abc.transaction.counterpary.app.com",
    "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
    "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
    "timestamp": "2022-10-28T10:34:58.469Z"
  },
  "message":
  {
    "issue_resolution_details":
    {
      "issue":
      {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "complainant_info":
        {
          "person":
          {
            "name": "Sam Manuel"
          },
          "contact":
          {
            "name": "Sam Manuel",
            "address":
            {
              "full": "45,Park Avenue, Delhi",
              "format": "string"
            },
            "phone": "9879879870",
            "email": "sam@yahoo.com"
          },
          "complainant_action": "open",
          "complainant_action_reason": "identified an issue with order so opening an issue"
        },
        "order":
        {
          "id": "0949848944",
          "ref_order_ids":
          [
            "1209",
            "1210"
          ],
          "items":
          [ {
              "id": "23433",
              "parent_item_id": "44334",
              "descriptor":
              {
                "name": "laptop",
                "code": "code",
                "short_desc": "laptop description",
                "long_desc": "laptop long description",
                "additional_desc":
                {
                  "url": "string",
                  "content_type": "text/plain"
                }
              },
              "manufacturer":
              {
                "descriptor":
                {
                  "name": "manuf desc",
                  "code": "mfg93",
                  "short_desc": "mfg desc ",
                  "long_desc": "mfg long desc"
                },
                "address":
                {
                  "full": "433054, Park Avenue, Delhi",
                  "format": "string"
                },
                "contact":
                {
                  "name": "Sam X",
                  "address":
                  {
                    "full": "4993, Park Avenue ,Delhi",
                    "format": "string"
                  },
                  "phone": "4059404549",
                  "email": "email@email.com"
                }
              },
              "price":
              {
                "currency": "INR",
                "value": "590"
              },
              "quantity":
              {
                "allocated":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "3",
                    "unit": "string"
                  }
                },
                "available":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "6",
                    "unit": "string"
                  }
                },
                "maximum":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "2",
                    "unit": "string"
                  }
                }
              },
              "category_id": "409",
              "fulfillment_id": "230"
            }
          ],
          "fulfillments":
          [ {
              "id": "230",
              "type": "xyz",
              "provider_id": "30493",
              "state":
              {
                "descriptor":
                {
                  "name": "fulfillments state name",
                  "code": "fulfillments state code",
                  "short_desc": "fulfillments desc",
                  "long_desc": "fulfillments long desc name"
                },
                "updated_at": "2022-12-08T05:51:34.946Z",
                "updated_by": "string"
              },
              "tracking": false,
              "customer":
              {
                "person":
                {
                  "id": "101",
                  "name": "Sam B"
                },
                "contact":
                {
                  "name": "Sam B",
                  "address":
                  {
                    "full": "43450, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "5904500590",
                  "email": "sam@emial.com"
                }
              },
              "agent":
              {
                "person":
                {
                  "id": "40904",
                  "name": "Sam A"
                },
                "contact":
                {
                  "name": "Sam F C",
                  "address":
                  {
                    "full": "54059, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "4304903440",
                  "email": "email@pa.com"
                },
                "organization":
                {
                  "descriptor":
                  {
                    "name": "FAOrg",
                    "code": "09403",
                    "short_desc": "FA ORG shot desc",
                    "long_desc": "FA Logn shot desc"
                  },
                  "address":
                  {
                    "full": "494589, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "contact":
                  {
                    "name": "Sam FAO",
                    "address":
                    {
                      "full": "95904, Park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "9403904940",
                    "email": "fao@email.com"
                  }
                }
              },
              "contact":
              {
                "name": "con",
                "address":
                {
                  "full": "4390, Park Avenue Delhi",
                  "format": "string"
                },
                "phone": "0549505400",
                "email": "dmai@email.com"
              }
            }
          ],
          "created_at": "2022-12-08T05:51:34.961Z",
          "updated_at": "2022-12-08T05:51:34.961Z"
        },
        "description":
        {
          "name": "issue is with the product",
          "short_desc": "product is duplicate",
          "long_desc": "the product appears to be of bad quality",
          "additional_desc":
          {
            "url": "https://link-to-assests.buyer.com/product-imag.img",
            "content_type": "text/plain"
          }
        },
        "category": "Product",
        "issue_type": "Issue",
        "source":
        {
          "network_participant_id": "813389c0-65a2-11ed-9022-0242ac120002",
          "issue_source_type": "Consumer"
        },
        "supplementary_information":
        [ {
            "issue_update_info":
            {
              "name": "string",
              "code": "respondent",
              "short_desc": "provide further proofs of the issue",
              "long_desc": "provide some pictures to verify if the product is actually not of the right quality"
            },
            "created_at": "2022-12-04T06:30:47.753Z",
            "modified_at": "2022-12-04T06:30:47.753Z"
          }
        ],
        "expected_response_time":
        {
          "duration": "T1H"
        },
        "expected_resolution_time":
        {
          "duration": "P1D"
        },
        "status":
        {
          "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "status": "Resolved",
          "status_details":
          {
            "code": "Accepted",
            "short_desc": "issue status details ",
            "long_desc": "issue status long description"
          },
          "status_change_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "issue_modification_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "modified_by":
          {
            "org":
            {
              "name": "org name of status modifier"
            },
            "contact":
            {
              "name": "contact name of status modifier",
              "address":
              {
                "full": "40595, park Avenue Delhi",
                "format": "string"
              },
              "phone": "04950394039",
              "email": "modiemail@org.com"
            },
            "person":
            {
              "name": "person name of status modifier"
            }
          },
          "modified_issue_field_name": "string"
        },
        "created_at": "2022-12-04T06:30:47.753Z",
        "modified_at": "2022-12-04T06:30:47.753Z"
      },
      "resolution_provider":
      {
        "respondent_info":
        {
          "type": "Transaction Counterparty NP",
          "organization":
          {
            "org":
            {
              "name": "interfacing app resolution provider app name"
            },
            "contact":
            {
              "name": "Sam C",
              "address":
              {
                "full": "4959, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "4059304940",
              "email": "email@resolutionproviderorg.com"
            },
            "person":
            {
              "name": "resolution provider org contact person name"
            }
          },
          "resolution_support":
          {
            "respondentChatLink": "http://chat-link/respondant",
            "respondentEmail": "email@respondant.com",
            "respondentContact":
            {
              "name": "resolution providers, issue resolution support respondent name",
              "address":
              {
                "full": "4955, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "5949595059",
              "email": "respondantemail@resolutionprovider.com"
            },
            "respondentFaqs":
            {
              "additionalProp1":
              {
                "question": "Question 1",
                "Answer": "Answer 1"
              },
              "additionalProp2":
              {
                "question": "question 2",
                "Answer": "answer 2"
              },
              "additionalProp3":
              {
                "question": "Question 3",
                "Answer": "Answer 3"
              }
            },
            "additional_sources":
            {
              "additionalProp1":
              {
                "type": "additional resolution source type : eg : IVR",
                "link": "ivr contact details or link to it : eg: http://resolutionsupport.ivr.resolutionprovider.com"
              }
            },
            "gros":
            [ {
              "person": {
                "name": "Sam D",
                "image": "https://gro/profile-image.png"
              },
              "contact":
              {
                "name": "Sam D",
                "address":
                {
                  "full": "45 park anv. Delhi",
                  "format": "string"
                },
                "phone": "9605960796",
                "email": "email@gro.com"
              },
              "gro_type": "Interfacing NP GRO"
              }
          ],
            "selected_odrs":
            [ {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization":
                {
                  "org":
                  {
                    "name": "ODR provider"
                  },
                  "contact":
                  {
                    "name": "ODR provider contact name",
                    "address":
                    {
                      "full": "A-3040, park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "0495049304",
                    "email": "email@odr.com"
                  },
                  "person":
                  {
                    "name": "odr provider contact  person name"
                  }
                },
                "pricing_model":
                {
                  "price":
                   {
                    "currency": "rupees",
                    "value": "10000"
                   },
                  "pricing_info": "10% of invoice cost"
                },
                "resolution_ratings":
                {
                  "rating_value": "98%"
                }
              }
            ]
          }
        }
      },
      "resolution": {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "resolution": "issue resolution details",
        "resolution_remarks": "remarks related to the resolution",
        "gro_remarks": "remarks from gro related to issue resolution",
        "dispute_resolution_remarks": "provided by odr about dispute resolution",
        "resolution_action": "resolve"
      }
    }
  },
  "error":
  {
    "code": "string",
    "path": "string",
    "message": "string"
  }
}');
INSERT INTO public.reports(id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at,report_data)VALUES(15, true, true, 1, '2022-06-21 18:27:06.065309+05:30', 1, '2022-06-21 18:27:06.065309+05:30', 1, NULL,
'{
  "context":
  {
    "domain": "nic2004:52110",
    "country": "IND",
    "city": "std:080",
    "action": "on_issue",
    "core_version": "1.0.0",
    "bap_id": "ondc.bap.com",
    "bap_uri": "https://abc.interfacing.app.com",
    "bpp_id": "ondc.bpp.com",
    "bpp_uri": "https://abc.transaction.counterpary.app.com",
    "transaction_id": "6baa811a-6cbe-4ad3-94e9-cbf96aaff343",
    "message_id": "5abb922b-6cbe-4ad3-94e9-cbf96bbff343",
    "timestamp": "2022-10-28T10:34:58.469Z"
  },
  "message":
  {
    "issue_resolution_details":
    {
      "issue":
      {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "complainant_info":
        {
          "person":
          {
            "name": "Sam Manuel"
          },
          "contact":
          {
            "name": "Sam Manuel",
            "address":
            {
              "full": "45,Park Avenue, Delhi",
              "format": "string"
            },
            "phone": "9879879870",
            "email": "sam@yahoo.com"
          },
          "complainant_action": "open",
          "complainant_action_reason": "identified an issue with order so opening an issue"
        },
        "order":
        {
          "id": "0949848944",
          "ref_order_ids":
          [
            "1209",
            "1210"
          ],
          "items":
          [ {
              "id": "23433",
              "parent_item_id": "44334",
              "descriptor":
              {
                "name": "laptop",
                "code": "code",
                "short_desc": "laptop description",
                "long_desc": "laptop long description",
                "additional_desc":
                {
                  "url": "string",
                  "content_type": "text/plain"
                }
              },
              "manufacturer":
              {
                "descriptor":
                {
                  "name": "manuf desc",
                  "code": "mfg93",
                  "short_desc": "mfg desc ",
                  "long_desc": "mfg long desc"
                },
                "address":
                {
                  "full": "433054, Park Avenue, Delhi",
                  "format": "string"
                },
                "contact":
                {
                  "name": "Sam X",
                  "address":
                  {
                    "full": "4993, Park Avenue ,Delhi",
                    "format": "string"
                  },
                  "phone": "4059404549",
                  "email": "email@email.com"
                }
              },
              "price":
              {
                "currency": "INR",
                "value": "590"
              },
              "quantity":
              {
                "allocated":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "3",
                    "unit": "string"
                  }
                },
                "available":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "6",
                    "unit": "string"
                  }
                },
                "maximum":
                {
                  "count": 1,
                  "measure":
                  {
                    "type": "CONSTANT",
                    "value": "2",
                    "unit": "string"
                  }
                }
              },
              "category_id": "409",
              "fulfillment_id": "230"
            }
          ],
          "fulfillments":
          [ {
              "id": "230",
              "type": "xyz",
              "provider_id": "30493",
              "state":
              {
                "descriptor":
                {
                  "name": "fulfillments state name",
                  "code": "fulfillments state code",
                  "short_desc": "fulfillments desc",
                  "long_desc": "fulfillments long desc name"
                },
                "updated_at": "2022-12-08T05:51:34.946Z",
                "updated_by": "string"
              },
              "tracking": false,
              "customer":
              {
                "person":
                {
                  "id": "101",
                  "name": "Sam B"
                },
                "contact":
                {
                  "name": "Sam B",
                  "address":
                  {
                    "full": "43450, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "5904500590",
                  "email": "sam@emial.com"
                }
              },
              "agent":
              {
                "person":
                {
                  "id": "40904",
                  "name": "Sam A"
                },
                "contact":
                {
                  "name": "Sam F C",
                  "address":
                  {
                    "full": "54059, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "phone": "4304903440",
                  "email": "email@pa.com"
                },
                "organization":
                {
                  "descriptor":
                  {
                    "name": "FAOrg",
                    "code": "09403",
                    "short_desc": "FA ORG shot desc",
                    "long_desc": "FA Logn shot desc"
                  },
                  "address":
                  {
                    "full": "494589, Park Avenue, Delhi",
                    "format": "string"
                  },
                  "contact":
                  {
                    "name": "Sam FAO",
                    "address":
                    {
                      "full": "95904, Park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "9403904940",
                    "email": "fao@email.com"
                  }
                }
              },
              "contact":
              {
                "name": "con",
                "address":
                {
                  "full": "4390, Park Avenue Delhi",
                  "format": "string"
                },
                "phone": "0549505400",
                "email": "dmai@email.com"
              }
            }
          ],
          "created_at": "2022-12-08T05:51:34.961Z",
          "updated_at": "2022-12-08T05:51:34.961Z"
        },
        "description":
        {
          "name": "issue is with the product",
          "short_desc": "product is duplicate",
          "long_desc": "the product appears to be of bad quality",
          "additional_desc":
          {
            "url": "https://link-to-assests.buyer.com/product-imag.img",
            "content_type": "text/plain"
          }
        },
        "category": "Product",
        "issue_type": "Issue",
        "source":
        {
          "network_participant_id": "813389c0-65a2-11ed-9022-0242ac120002",
          "issue_source_type": "Consumer"
        },
        "supplementary_information":
        [ {
            "issue_update_info":
            {
              "name": "string",
              "code": "respondent",
              "short_desc": "provide further proofs of the issue",
              "long_desc": "provide some pictures to verify if the product is actually not of the right quality"
            },
            "created_at": "2022-12-04T06:30:47.753Z",
            "modified_at": "2022-12-04T06:30:47.753Z"
          }
        ],
        "expected_response_time":
        {
          "duration": "T1H"
        },
        "expected_resolution_time":
        {
          "duration": "P1D"
        },
        "status":
        {
          "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "status": "Resolved",
          "status_details":
          {
            "code": "Accepted",
            "short_desc": "issue status details ",
            "long_desc": "issue status long description"
          },
          "status_change_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "issue_modification_date":
          {
            "timestamp": "2022-12-04T06:30:47.753Z"
          },
          "modified_by":
          {
            "org":
            {
              "name": "org name of status modifier"
            },
            "contact":
            {
              "name": "contact name of status modifier",
              "address":
              {
                "full": "40595, park Avenue Delhi",
                "format": "string"
              },
              "phone": "04950394039",
              "email": "modiemail@org.com"
            },
            "person":
            {
              "name": "person name of status modifier"
            }
          },
          "modified_issue_field_name": "string"
        },
        "created_at": "2022-12-04T06:30:47.753Z",
        "modified_at": "2022-12-04T06:30:47.753Z"
      },
      "resolution_provider":
      {
        "respondent_info":
        {
          "type": "Transaction Counterparty NP",
          "organization":
          {
            "org":
            {
              "name": "interfacing app resolution provider app name"
            },
            "contact":
            {
              "name": "Sam C",
              "address":
              {
                "full": "4959, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "4059304940",
              "email": "email@resolutionproviderorg.com"
            },
            "person":
            {
              "name": "resolution provider org contact person name"
            }
          },
          "resolution_support":
          {
            "respondentChatLink": "http://chat-link/respondant",
            "respondentEmail": "email@respondant.com",
            "respondentContact":
            {
              "name": "resolution providers, issue resolution support respondent name",
              "address":
              {
                "full": "4955, Park Avenue Delhi",
                "format": "string"
              },
              "phone": "5949595059",
              "email": "respondantemail@resolutionprovider.com"
            },
            "respondentFaqs":
            {
              "additionalProp1":
              {
                "question": "Question 1",
                "Answer": "Answer 1"
              },
              "additionalProp2":
              {
                "question": "question 2",
                "Answer": "answer 2"
              },
              "additionalProp3":
              {
                "question": "Question 3",
                "Answer": "Answer 3"
              }
            },
            "additional_sources":
            {
              "additionalProp1":
              {
                "type": "additional resolution source type : eg : IVR",
                "link": "ivr contact details or link to it : eg: http://resolutionsupport.ivr.resolutionprovider.com"
              }
            },
            "gros":
            [ {
              "person": {
                "name": "Sam D",
                "image": "https://gro/profile-image.png"
              },
              "contact":
              {
                "name": "Sam D",
                "address":
                {
                  "full": "45 park anv. Delhi",
                  "format": "string"
                },
                "phone": "9605960796",
                "email": "email@gro.com"
              },
              "gro_type": "Interfacing NP GRO"
              }
          ],
            "selected_odrs":
            [ {
                "name": "ODR provider name",
                "about_info": "information about the ODR provider",
                "url": "https://info.odrprovider.com/",
                "organization":
                {
                  "org":
                  {
                    "name": "ODR provider"
                  },
                  "contact":
                  {
                    "name": "ODR provider contact name",
                    "address":
                    {
                      "full": "A-3040, park Avenue Delhi",
                      "format": "string"
                    },
                    "phone": "0495049304",
                    "email": "email@odr.com"
                  },
                  "person":
                  {
                    "name": "odr provider contact  person name"
                  }
                },
                "pricing_model":
                {
                  "price":
                   {
                    "currency": "rupees",
                    "value": "10000"
                   },
                  "pricing_info": "10% of invoice cost"
                },
                "resolution_ratings":
                {
                  "rating_value": "98%"
                }
              }
            ]
          }
        }
      },
      "resolution": {
        "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "resolution": "issue resolution details",
        "resolution_remarks": "remarks related to the resolution",
        "gro_remarks": "remarks from gro related to issue resolution",
        "dispute_resolution_remarks": "provided by odr about dispute resolution",
        "resolution_action": "resolve"
      }
    }
  },
  "error":
  {
    "code": "string",
    "path": "string",
    "message": "string"
  }
}');
SELECT setval( pg_get_serial_sequence('public.reports', 'id'), (SELECT max(id) FROM public.reports));




SET session_replication_role = 'origin';























