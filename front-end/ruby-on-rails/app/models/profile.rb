class Profile < ApplicationRecord
	belongs_to :experiment
	has_one :cell, through: :experiment
	has_one :drug, through: :experiment
	has_one :dataset, through: :experiment
end
