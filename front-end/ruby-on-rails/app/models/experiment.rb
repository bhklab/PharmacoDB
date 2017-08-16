class Experiment < ApplicationRecord
	has_many :dose_responses
	has_one :profile
	belongs_to :cell
	belongs_to :drug
	belongs_to :dataset
	belongs_to :tissue
	has_many :source_cell_names, through: :cell
	has_many :source_drug_names, through: :drug
	has_many :drug_targets, through: :drug
	has_many :targets, through: :drug_targets
end
