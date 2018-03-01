class AutocompleteController < ApplicationController

    #### TODO:: when printing out the autocomplete results, we can order categories by the # of 
    #    terms used to match them, and then loop from max to min printing out matches until we hit maximum print # (15?)

    def suggest

        input = params.values[0].upcase.split(Rails.application.config.my_search_split)
        
        if !params[:query].present? & !params[:query].nil?
            render json: []
        end

        single_input = input.length.downto(1).flat_map{|size| input.combination(size).map{|x| [x.join(" "), x.length]}}
        cell_result = single_input.map{|term| CELLS.children_with_values(term[0]).map{|x| [Cell.find(x[1]).cell_name, term[1]]}}.flatten(1).sort_by{|x| -x[1]}.uniq{|x| x[0]}
        tissue_result = single_input.map{|term| TISSUES.children_with_values(term[0]).map{|x| [Tissue.find(x[1]).tissue_name, term[1]]}}.flatten(1).sort_by{|x| -x[1]}.uniq{|x| x[0]}
        drug_result =  single_input.map{|term| DRUGS.children_with_values(term[0]).map{|x| [Drug.find(x[1]).drug_name, term[1]]}}.flatten(1).sort_by{|x| -x[1]}.uniq{|x| x[0]}
        gene_result = single_input.map{|term| GENES.children_with_values(term[0]).map{|x| [Gene.find(x[1]).gene_name, term[1]]}}.flatten(1).sort_by{|x| -x[1]}.uniq{|x| x[0]}
        dataset_result = single_input.map{|term| DATASETS.children_with_values(term[0]).map{|x| [Dataset.find(x[1]).dataset_name, term[1]]}}.flatten(1).sort_by{|x| -x[1]}.uniq{|x| x[0]}

        ### TODO filter by experiments present!

        cell_drug_result = []
        unless input.length < 2
            (input.length-1).downto(1).each do |i|
                first = input[0...(i)]
                second = input[i...(input.length)]

                temp = first.length.downto(1).flat_map{|size| first.combination(size).map{|x| [x.join(" "), x.length]}}
                first_matches_cell = temp.flat_map{|term| CELLS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                first_matches_drug = temp.flat_map{|term| DRUGS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                
                temp = second.length.downto(1).flat_map{|size| second.combination(size).map{|x| [x.join(" "), x.length]}}
                second_matches_cell = temp.flat_map{|term| CELLS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                second_matches_drug = temp.flat_map{|term| DRUGS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                
                matches_one = first_matches_cell.product(second_matches_drug)
                matches_two = first_matches_drug.product(second_matches_cell)
                matches_one = matches_one.select{|x| !Experiment.where(cell_id: x[0][0], drug_id: x[1][0]).empty?}.map{|x| [Cell.find(x[0][0]).cell_name + " " + Drug.find(x[1][0]).drug_name, x[0][1] + x[1][1]]}
                matches_two = matches_two.select{|x| !Experiment.where(drug_id: x[0][0], cell_id: x[1][0]).empty?}.map{|x| [Drug.find(x[0][0]).drug_name + " " + Cell.find(x[1][0]).cell_name, x[0][1] + x[1][1]]}
                cell_drug_result += matches_one + matches_two


            end
        end

        ### XXX:: FIXME:: CURRENTLY A HACK as it works only when dataset is 1 term
        ### TODO:: transfer changes from above
        unless input.length < 3
            third = input[input.length-1]


            dataset_matches = DATASETS.children_with_values(third).map{|x| [x[1], 1]}
            unless dataset_matches.empty?

                (input.length-2).downto(1).each do |i|
                    first = input[0...(i)]
                    second = input[i...(input.length-1)]

                    temp = first.length.downto(1).flat_map{|size| first.combination(size).map{|x| [x.join(" "), x.length]}}
                    first_matches_cell = temp.flat_map{|term| CELLS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                    first_matches_drug = temp.flat_map{|term| DRUGS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                    
                    temp = second.length.downto(1).flat_map{|size| second.combination(size).map{|x| [x.join(" "), x.length]}}
                    second_matches_cell = temp.flat_map{|term| CELLS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                    second_matches_drug = temp.flat_map{|term| DRUGS.children_with_values(term[0]).map{|x| [x[1], term[1]]}}#.sort_by{|x| -x[1]}.select{|x| !x[0].empty?}
                                     
                    matches_one = first_matches_cell.product(second_matches_drug, dataset_matches)
                    matches_two = first_matches_drug.product(second_matches_cell, dataset_matches)

                    matches_one = matches_one.select{|x| !Experiment.where(cell_id: x[0][0], drug_id: x[1][0], dataset_id: x[2][0]).empty?}
                    matches_two = matches_two.select{|x| !Experiment.where(drug_id: x[0][0], cell_id: x[1][0], dataset_id: x[2][0]).empty?}

                    matches_one = matches_one.map{|x| [Cell.find(x[0][0]).cell_name + " " + Drug.find(x[1][0]).drug_name + " " + Dataset.find(x[2][0]).dataset_name, x[0][1] + x[1][1] + x[2][1]]}
                    matches_two = matches_two.map{|x| [Drug.find(x[0][0]).drug_name + " " + Cell.find(x[1][0]).cell_name + " " + Dataset.find(x[2][0]).dataset_name, x[0][1] + x[1][1] + x[2][1]]}
                    
                    cell_drug_result += matches_one + matches_two + cell_drug_result

                end
            end

        end


        cell_drug_result = cell_drug_result.sort_by{|x| -x[1]}.uniq{|x| x[0]}

        ## XXX:: Using fact that datasets are 1 word only, optimization/hack 

        two_dataset_result = input.combination(2).
                                   map{|x| x.map{|term| DATASETS.children_with_values(term).map{|y| Dataset.find(y[1]).dataset_name}}}.
                                   flat_map{|x| x[0].product(x[1])}.
                                   select{|x| x[0] != x[1]}.map{|x| [x.join(" "),2]}

        three_dataset_result = input.combination(3).
                                     map{|x| x.map{|term| DATASETS.children_with_values(term).map{|y| Dataset.find(y[1]).dataset_name}}}.
                                     flat_map{|x| x[0].product(x[1], x[2]).select{|x| x.length == x.uniq.length}.
                                     map{|y| [y.join(" "),3]}}

        dataset_int_result = three_dataset_result + two_dataset_result

        ### TODO:: dataset intersection

        # input = (input.length-1).downto(1).map{|size| input.combination(size).map{|x| x.join(" ")}}
        render json: {:cell => cell_result, 
                      :tissue => tissue_result, 
                      :drug => drug_result, 
                      :gene => gene_result, 
                      :dataset => dataset_result, 
                      :cell_drug => cell_drug_result,
                      :dataset_int => dataset_int_result}.to_json
    end
end